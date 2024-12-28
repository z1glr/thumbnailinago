package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ClientSettings struct {
	ReplacementPattern string     `json:"replacement_pattern"`
	DateFormat         string     `json:"date_format"`
	ExportDays         ExportDays `json:"export_days"`
}

type config struct {
	LastSvg        string         `json:"last_svg"`
	LastExport     string         `json:"last_export"`
	ClientSettings ClientSettings `json:"client_settings"`
}

var Config config

func init() {
	if configStr, err := os.ReadFile("config.json"); err != nil {
		panic(fmt.Errorf("can't read config-file: %v", err))
	} else {
		if err := json.Unmarshal(configStr, &Config); err != nil {
			panic(fmt.Errorf("can't parse config-file: %v", err))
		}

		// try to load the last opened svg
		if f, err := os.ReadFile(Config.LastSvg); err == nil {
			svg.Name = Config.LastSvg

			svg.Str = string(f)
		}
	}
}

func writeConfig() {
	if strConfig, err := json.MarshalIndent(Config, "", "\t"); err != nil {
		panic(fmt.Errorf("can't stringify config: %v", err))
	} else {
		os.WriteFile("config.json", strConfig, 0o644)
	}
}
