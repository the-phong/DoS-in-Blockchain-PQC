package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"example.com/pluto/internal/node"
	"example.com/pluto/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stdout, "Pluto Phase 02: initialize a local node, then start it.")
		usage(stdout)
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "help takes no arguments")
			return 2
		}
		usage(stdout)
		return 0
	case "version", "-v", "--version":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "version takes no arguments")
			return 2
		}
		fmt.Fprintf(stdout, "plutod %s\nGo: %s\nPlatform: %s/%s\nRevision: %s\n", version.Number, runtime.Version(), runtime.GOOS, runtime.GOARCH, version.Revision)
		return 0
	case "init":
		flags := flag.NewFlagSet("init", flag.ContinueOnError)
		flags.SetOutput(stderr)
		home := flags.String("home", ".pluto", "node data directory")
		chainID := flags.String("chain-id", "pluto-dev-1", "chain identity")
		rpcPort := flags.Int("rpc-port", 26657, "localhost RPC port")
		p2pPort := flags.Int("p2p-port", 26656, "localhost P2P port")
		abciPort := flags.Int("abci-port", 26658, "localhost ABCI port")
		if err := flags.Parse(args[1:]); err != nil || flags.NArg() != 0 {
			fmt.Fprintln(stderr, "invalid init arguments")
			return 2
		}
		if err := node.Init(*home, cometBinary(), *chainID, *rpcPort, *p2pPort, *abciPort, stdout); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Initialized Pluto node at %s (chain %s)\n", *home, *chainID)
		return 0
	case "start":
		flags := flag.NewFlagSet("start", flag.ContinueOnError)
		flags.SetOutput(stderr)
		home := flags.String("home", ".pluto", "node data directory")
		if err := flags.Parse(args[1:]); err != nil || flags.NArg() != 0 {
			fmt.Fprintln(stderr, "invalid start arguments")
			return 2
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if err := node.Run(ctx, *home, cometBinary(), stdout); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	case "state":
		flags := flag.NewFlagSet("state", flag.ContinueOnError)
		flags.SetOutput(stderr)
		home := flags.String("home", ".pluto", "node data directory")
		if err := flags.Parse(args[1:]); err != nil || flags.NArg() != 0 {
			fmt.Fprintln(stderr, "invalid state arguments")
			return 2
		}
		state, err := node.State(*home)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		result := struct {
			ChainID string `json:"chain_id"`
			Version uint32 `json:"state_version"`
			Height  int64  `json:"height"`
			AppHash string `json:"app_hash"`
		}{state.ChainID, state.Version, state.Height, hex.EncodeToString(state.AppHash[:])}
		if err := json.NewEncoder(stdout).Encode(result); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		usage(stderr)
		return 2
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: plutod [init|start|state|help|version] [flags]")
	fmt.Fprintln(w, "Phase 02 supports empty blocks only; application transactions are rejected.")
}

func cometBinary() string {
	self, err := os.Executable()
	if err != nil {
		return "cometbft-unavailable"
	}
	return filepath.Join(filepath.Dir(self), "cometbft")
}
