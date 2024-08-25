package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigurationItf interface {
	LoadConfiguration(path string) error
	GetConfiguration() Configuration
}

type Configuration struct {
	Server   Server
	Resource Resource
}

func NewConfiguration(path string) ConfigurationItf {
	config := Configuration{}
	config.LoadConfiguration(path)
	return &config
}

// LoadConfiguration loads the configuration from the provided file path
// and fills in the configuration struct. In production, proper security
// mechanisms should be implemented to protect configuration data.
func (cfg *Configuration) LoadConfiguration(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	byteFile, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteFile, cfg)
	if err != nil {
		return err
	}

	return nil
}

// GetConfiguration returns the current configuration.
func (cfg *Configuration) GetConfiguration() Configuration {
	return *cfg
}
