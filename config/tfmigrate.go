package config

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/minamijoyo/tfmigrate/history"
)

// ConfigurationFile represents a file for CLI settings in HCL.
type ConfigurationFile struct {
	// Tfmigrate is a top-level block.
	// It must contain only one block, and multiple blocks are not allowed.
	Tfmigrate TfmigrateBlock `hcl:"tfmigrate,block"`
}

// TfmigrateBlock represents a block for CLI settings in HCL.
type TfmigrateBlock struct {
	// MigrationDir is a path to directory where migration files are stored.
	// Default to `.` (current directory).
	MigrationDir string `hcl:"migration_dir,optional"`
	// History is a block for migration history management.
	History *HistoryBlock `hcl:"history,block"`
}

// TfmigrateConfig is a config for top-level CLI settings.
// TfmigrateBlock is just used for parsing HCL and
// TfmigrateConfig is used for building application logic.
// TODO: Here is where config is defined/used
type TfmigrateConfig struct {
	// MigrationDir is a path to directory where migration files are stored.
	// Default to `.` (current directory).
	MigrationDir string
	// History is a config for migration history management.
	History *history.Config
}

// LoadConfigurationFile is a helper function which reads and parses a given configuration file.
func LoadConfigurationFile(filename string) (*TfmigrateConfig, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseConfigurationFile(filename, source)
}

// ParseConfigurationFile parses a given source of configuration file and
// returns a TfmigrateConfig.
// Note that this method does not read a file and you should pass source of config in bytes.
// The filename is used for error message and selecting HCL syntax (.hcl and .json).
// TODO: Here is where the actual configuration is loaded.
func ParseConfigurationFile(filename string, source []byte) (*TfmigrateConfig, error) {
	// Decode tfmigrate block.
	var f ConfigurationFile
	err := hclsimple.Decode(filename, source, nil, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode setting file: %s, err: %s", filename, err)
	}

	config := NewDefaultConfig()
	if len(f.Tfmigrate.MigrationDir) > 0 {
		config.MigrationDir = f.Tfmigrate.MigrationDir
	}

	if f.Tfmigrate.History != nil {
		history, err := parseHistoryBlock(*f.Tfmigrate.History)
		if err != nil {
			return nil, err
		}
		config.History = history
	}

	return config, nil
}

// NewDefaultConfig returns a new instance of TfmigrateConfig.
func NewDefaultConfig() *TfmigrateConfig {
	return &TfmigrateConfig{
		MigrationDir: ".",
	}
}
