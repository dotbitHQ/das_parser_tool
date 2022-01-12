package config

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/scorpiotzh/mylog"
	"github.com/scorpiotzh/toolib"
)

var (
	Cfg CfgServer
	log = mylog.NewLogger("config", mylog.LevelDebug)
)

func InitCfg(configFilePath string) error {
	log.Info("read from config：", configFilePath)
	if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
		return fmt.Errorf("UnmarshalYamlFile err:%s", err.Error())
	}
	log.Info("config file：", toolib.JsonString(Cfg))
	return nil
}

type CfgServer struct {
	Chain struct {
		Net      common.DasNetType `json:"net" yaml:"net"`
		CkbUrl   string            `json:"ckb_url" yaml:"ckb_url"`
		IndexUrl string            `json:"index_url" yaml:"index_url"`
	} `json:"chain" yaml:"chain"`
	DasCore struct {
		THQCodeHash         string                            `json:"thq_code_hash" yaml:"thq_code_hash"`
		DasContractArgs     string                            `json:"das_contract_args" yaml:"das_contract_args"`
		DasContractCodeHash string                            `json:"das_contract_code_hash" yaml:"das_contract_code_hash"`
		MapDasContract      map[common.DasContractName]string `json:"map_das_contract" yaml:"map_das_contract"`
	} `json:"das_core" yaml:"das_core"`
}
