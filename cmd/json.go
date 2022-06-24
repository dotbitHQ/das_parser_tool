package cmd

import (
	"context"
	"das_parser_tool/chain"
	"das_parser_tool/config"
	"das_parser_tool/dascore"
	"das_parser_tool/parser"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var (
	jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Parser transaction by transaction json",
	}
	jsonFileCmd = &cobra.Command{
		Use:   "file",
		Short: "Parser transaction by transaction json file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file, err := ioutil.ReadFile(args[0])
			if err != nil {
				cobra.CheckErr(fmt.Errorf("ReadFile err: %v ", err.Error()))
			}

			jsonParser(file)
		},
	}
	jsonDataCmd = &cobra.Command{
		Use:   "data",
		Short: "Parser transaction by transaction json data",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, v := range args {
				jsonParser([]byte(v))
			}
		},
	}
)

func init() {
	jsonCmd.AddCommand(jsonFileCmd)
	jsonCmd.AddCommand(jsonDataCmd)
}

type Transaction struct {
	CellDeps []struct {
		DepType  string `json:"dep_type"`
		OutPoint struct {
			Index  string `json:"index"`
			TxHash string `json:"tx_hash"`
		} `json:"out_point"`
	} `json:"cell_deps"`
	Hash       string        `json:"hash"`
	HeaderDeps []interface{} `json:"header_deps"`
	Inputs     []struct {
		PreviousOutput struct {
			Index  string `json:"index"`
			TxHash string `json:"tx_hash"`
		} `json:"previous_output"`
		Since string `json:"since"`
	} `json:"inputs"`
	Outputs []struct {
		Capacity string `json:"capacity"`
		Lock     struct {
			Args     string `json:"args"`
			CodeHash string `json:"code_hash"`
			HashType string `json:"hash_type"`
		} `json:"lock"`
		Type *struct {
			Args     string `json:"args"`
			CodeHash string `json:"code_hash"`
			HashType string `json:"hash_type"`
		} `json:"type"`
	} `json:"outputs"`
	OutputsData []string `json:"outputs_data"`
	Version     string   `json:"version"`
	Witnesses   []string `json:"witnesses"`
}

func jsonParser(arg []byte) {
	var tx Transaction
	err := json.Unmarshal(arg, &tx)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Unmarshal err: %v ", err.Error()))
	}

	if tx.Hash != "" {
		hashParser(tx.Hash)
		return
	}

	// ckb node
	ckbClient := chain.NewClient(context.Background(), config.Cfg.Chain.CkbUrl, config.Cfg.Chain.IndexUrl)
	// contract init
	dasCore := dascore.NewDasCore(ckbClient.Client())
	// transaction parser
	bp := parser.Parser{
		CkbClient: ckbClient,
		DasCore:   dasCore,
	}
	transaction := convertTransaction(tx)
	out := bp.JsonParser(transaction)

	b, err := json.Marshal(out)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("Marshal err: %v ", err.Error()))
	}
	fmt.Println(string(b))
}

func convertTransaction(tx Transaction) *types.Transaction {
	version, _ := hexutil.DecodeUint64(tx.Version)
	var cellDeps []*types.CellDep
	for _, v := range tx.CellDeps {
		index, _ := hexutil.DecodeUint64(v.OutPoint.Index)
		cellDeps = append(cellDeps, &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash(v.OutPoint.TxHash),
				Index:  uint(index),
			},
			DepType: types.DepType(v.DepType),
		})
	}

	var inputs []*types.CellInput
	for _, v := range tx.Inputs {
		since, _ := hexutil.DecodeUint64(v.Since)
		index, _ := hexutil.DecodeUint64(v.PreviousOutput.Index)
		inputs = append(inputs, &types.CellInput{
			Since: since,
			PreviousOutput: &types.OutPoint{
				TxHash: types.HexToHash(v.PreviousOutput.TxHash),
				Index:  uint(index),
			},
		})
	}

	var outputs []*types.CellOutput
	var outputsData [][]byte
	for k, v := range tx.Outputs {
		outputsData = append(outputsData, common.Hex2Bytes(tx.OutputsData[k]))
		capacity, _ := hexutil.DecodeUint64(v.Capacity)
		switch v.Type {
		case nil:
			outputs = append(outputs, &types.CellOutput{
				Capacity: capacity,
				Lock: &types.Script{
					CodeHash: types.HexToHash(v.Lock.CodeHash),
					HashType: types.ScriptHashType(v.Lock.HashType),
					Args:     common.Hex2Bytes(v.Lock.Args),
				},
			})
		default:
			outputs = append(outputs, &types.CellOutput{
				Capacity: capacity,
				Lock: &types.Script{
					CodeHash: types.HexToHash(v.Lock.CodeHash),
					HashType: types.ScriptHashType(v.Lock.HashType),
					Args:     common.Hex2Bytes(v.Lock.Args),
				},
				Type: &types.Script{
					CodeHash: types.HexToHash(v.Type.CodeHash),
					HashType: types.ScriptHashType(v.Type.HashType),
					Args:     common.Hex2Bytes(v.Type.Args),
				},
			})
		}
	}

	var witnesses [][]byte
	for _, v := range tx.Witnesses {
		witnesses = append(witnesses, common.Hex2Bytes(v))
	}

	return &types.Transaction{
		Version:     uint(version),
		Hash:        types.HexToHash(tx.Hash),
		CellDeps:    cellDeps,
		Inputs:      inputs,
		Outputs:     outputs,
		OutputsData: outputsData,
		Witnesses:   witnesses,
	}
}
