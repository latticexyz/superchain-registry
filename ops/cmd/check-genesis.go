package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum-optimism/superchain-registry/ops/flags"
	"github.com/ethereum/go-ethereum/core"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/urfave/cli/v2"
)

var CheckGenesisCmd = cli.Command{
	Name:  "check-genesis",
	Flags: []cli.Flag{flags.GenesisFlag},
	Usage: "Sanity check genesis (genesis.json) is reproducible",
	Action: func(ctx *cli.Context) error {
		genesisPath := ctx.String(flags.GenesisFlag.Name)
		fmt.Printf("Attempting to read from %s\n", genesisPath)
		file, err := os.ReadFile(genesisPath)
		if err != nil {
			return fmt.Errorf("failed to read from local genesis.json config file: %w", err)
		}
		var localGenesis *core.Genesis
		if err = json.Unmarshal(file, &localGenesis); err != nil {
			return fmt.Errorf("failed to unmarshal local genesis.json into core.Genesis struct: %w", err)
		}

		chainId := localGenesis.Config.ChainID.Uint64()

		gethGenesis, err := core.LoadOPStackGenesis(chainId)
		if err != nil {
			return fmt.Errorf("failed to load genesis via op-geth: ensure chainId has already been added to registry: %w", err)
		}

		opts := cmp.Options{cmpopts.IgnoreUnexported(big.Int{})}
		if diff := cmp.Diff(localGenesis, gethGenesis, opts...); diff != "" {
			return fmt.Errorf("local genesis.json (-) does not match config generated by op-geth (+): %s", diff)
		}

		fmt.Println("👌 Regenerated genesis config matches existing one")
		return nil
	},
}
