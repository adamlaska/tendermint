package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

// InitFilesCmd initialises a fresh Tendermint Core instance.
var InitFilesCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Tendermint",
	RunE:  initFiles,
}

func initFiles(cmd *cobra.Command, args []string) error {
	return initFilesWithConfig(config)
}

func initFilesWithConfig(config *cfg.Config) error {
	// private validator
	privValkeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()
	var pv *privval.FilePV
	if cmn.FileExists(privValkeyFile) {
		pv = privval.LoadFilePV(privValkeyFile, privValStateFile)
		logger.Info("Found private validator", "keyFile", privValkeyFile,
			"stateFile", privValStateFile)
	} else {
		pv = privval.GenFilePV(privValkeyFile, privValStateFile)
		pv.Save()
		logger.Info("Generated private validator", "keyFile", privValkeyFile,
			"stateFile", privValStateFile)
	}

	nodeKeyFile := config.NodeKeyFile()
	if cmn.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := config.GenesisFile()
	if cmn.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc := types.GenesisDoc{
			ChainID:         fmt.Sprintf("test-chain-%v", cmn.RandStr(6)),
			GenesisTime:     tmtime.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}
		genDoc.Validators = []types.GenesisValidator{{
			Address: pv.GetPubKey().Address(),
			PubKey:  pv.GetPubKey(),
			Power:   10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}
