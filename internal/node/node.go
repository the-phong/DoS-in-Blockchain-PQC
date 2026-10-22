// Package node manages a local CometBFT process and Pluto ABCI socket server.
package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/cometbft/cometbft/abci/server"
	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/types"

	"example.com/pluto/internal/app"
	"example.com/pluto/internal/store"
)

type Config struct {
	ChainID     string `json:"chain_id"`
	ABCIAddress string `json:"abci_address"`
	RPCAddress  string `json:"rpc_address"`
	P2PAddress  string `json:"p2p_address"`
}

func address(port int) (string, error) {
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("invalid port %d", port)
	}
	return "tcp://127.0.0.1:" + strconv.Itoa(port), nil
}

func Init(home, cometBinary, chainID string, rpcPort, p2pPort, abciPort int, output io.Writer) error {
	if home == "" || cometBinary == "" {
		return errors.New("home and bundled CometBFT binary are required")
	}
	if _, err := store.Genesis(chainID); err != nil {
		return fmt.Errorf("invalid chain ID: %w", err)
	}
	if rpcPort == p2pPort || rpcPort == abciPort || p2pPort == abciPort {
		return errors.New("RPC, P2P, and ABCI ports must differ")
	}
	rpcAddr, err := address(rpcPort)
	if err != nil {
		return err
	}
	p2pAddr, err := address(p2pPort)
	if err != nil {
		return err
	}
	abciAddr, err := address(abciPort)
	if err != nil {
		return err
	}
	if _, err := os.Stat(home); err == nil {
		return fmt.Errorf("node home already exists: %s", home)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect node home: %w", err)
	}
	if err := os.MkdirAll(home, 0700); err != nil {
		return fmt.Errorf("create node home: %w", err)
	}
	cmd := exec.Command(cometBinary, "init", "--home", home)
	cmd.Stdout, cmd.Stderr = output, output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("initialize CometBFT files: %w", err)
	}
	genesisFile := filepath.Join(home, "config", "genesis.json")
	genesis, err := types.GenesisDocFromFile(genesisFile)
	if err != nil {
		return fmt.Errorf("read generated genesis: %w", err)
	}
	genesis.ChainID = chainID
	genesis.InitialHeight = 1
	genesis.AppState = nil
	if err := genesis.ValidateAndComplete(); err != nil {
		return fmt.Errorf("validate generated genesis: %w", err)
	}
	if err := genesis.SaveAs(genesisFile); err != nil {
		return fmt.Errorf("write generated genesis: %w", err)
	}
	config := cfg.DefaultConfig()
	config.SetRoot(home)
	config.ProxyApp = abciAddr
	config.RPC.ListenAddress = rpcAddr
	config.P2P.ListenAddress = p2pAddr
	config.Consensus.CreateEmptyBlocks = true
	config.Consensus.CreateEmptyBlocksInterval = time.Second
	cfg.WriteConfigFile(filepath.Join(home, "config", "config.toml"), config)
	metadata := Config{ChainID: chainID, ABCIAddress: abciAddr, RPCAddress: rpcAddr, P2PAddress: p2pAddr}
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(home, "config", "pluto.json"), append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("write Pluto node metadata: %w", err)
	}
	return nil
}

func ReadConfig(home string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(home, "config", "pluto.json"))
	if err != nil {
		return Config{}, fmt.Errorf("read Pluto node metadata: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("decode Pluto node metadata: %w", err)
	}
	if _, err := store.Genesis(config.ChainID); err != nil {
		return Config{}, fmt.Errorf("invalid Pluto node metadata: %w", err)
	}
	return config, nil
}

func State(home string) (*store.State, error) {
	if _, err := ReadConfig(home); err != nil {
		return nil, err
	}
	dbDir := filepath.Join(home, "data", "pluto")
	if _, err := os.Stat(dbDir); err != nil {
		return nil, fmt.Errorf("application database is absent: %w", err)
	}
	db, err := store.Open(dbDir)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	state, err := db.Load()
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, errors.New("application state has not been initialized")
	}
	return state, nil
}

// Run starts the Pluto socket server before launching the pinned CometBFT
// executable. Cancellation terminates only the child process started here.
func Run(ctx context.Context, home, cometBinary string, output io.Writer) error {
	config, err := ReadConfig(home)
	if err != nil {
		return err
	}
	genesis, err := types.GenesisDocFromFile(filepath.Join(home, "config", "genesis.json"))
	if err != nil {
		return fmt.Errorf("read CometBFT genesis: %w", err)
	}
	if genesis.ChainID != config.ChainID || genesis.InitialHeight != 1 || len(genesis.AppState) != 0 {
		return errors.New("CometBFT genesis conflicts with Pluto node metadata")
	}
	dbDir := filepath.Join(home, "data", "pluto")
	if _, err := os.Stat(dbDir); errors.Is(err, os.ErrNotExist) {
		if _, cometErr := os.Stat(filepath.Join(home, "data", "state.db")); cometErr == nil {
			return errors.New("Pluto database is absent while CometBFT state exists; refusing implicit reset")
		}
	} else if err != nil {
		return fmt.Errorf("inspect Pluto database: %w", err)
	}
	a, err := app.Open(dbDir)
	if err != nil {
		return err
	}
	defer a.Close()
	if a.Committed() == nil {
		if _, cometErr := os.Stat(filepath.Join(home, "data", "state.db")); cometErr == nil {
			return errors.New("Pluto committed state is absent while CometBFT state exists")
		}
	}
	if state := a.Committed(); state != nil && state.ChainID != config.ChainID {
		return errors.New("persisted application chain ID conflicts with genesis")
	}
	srv, err := server.NewServer(config.ABCIAddress, "socket", a)
	if err != nil {
		return err
	}
	if err := srv.Start(); err != nil {
		return fmt.Errorf("start ABCI server: %w", err)
	}
	defer srv.Stop()
	cmd := exec.Command(cometBinary, "node", "--home", home)
	cmd.Stdout, cmd.Stderr = output, output
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start CometBFT: %w", err)
	}
	result := make(chan error, 1)
	go func() { result <- cmd.Wait() }()
	select {
	case err := <-result:
		return fmt.Errorf("CometBFT exited: %w", err)
	case <-ctx.Done():
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("signal CometBFT: %w", err)
		}
		select {
		case err := <-result:
			if err != nil {
				return fmt.Errorf("stop CometBFT: %w", err)
			}
			return nil
		case <-time.After(10 * time.Second):
			_ = cmd.Process.Kill()
			<-result
			return errors.New("CometBFT did not stop within 10 seconds")
		}
	}
}
