package config

import (
	"log"

	"os"

	"path/filepath"

	"finala/serverutil" // Corrected import path based on module name

	"gopkg.in/yaml.v2" // Using existing v2 from go.mod
)

// AuthCredentialsConfig stores username and password.

type AuthCredentialsConfig struct {
	Username string `yaml:"username"`

	Password string `yaml:"password"`
}

// APIConfigYAML is the top-level structure for parsing api.yaml.

type APIConfigYAML struct {
	Auth AuthCredentialsConfig `yaml:"auth"`
}

// AppCredentials holds the effective username and password for login.

var AppCredentials AuthCredentialsConfig

// LoadCredentials reads auth credentials from the specified YAML file path.

func LoadCredentials(configPath string) {

	defaultUsername := "admin"

	configDir := filepath.Dir(configPath)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {

		log.Printf("INFO: Configuration directory %s does not exist. Creating.", configDir)

		if errMkdir := os.MkdirAll(configDir, 0755); errMkdir != nil {

			log.Fatalf("FATAL: Could not create configuration directory %s: %v", configDir, errMkdir)

		}

	}

	yamlFile, errReadFile := os.ReadFile(configPath)

	parsedConfig := APIConfigYAML{}

	configExistsAndReadable := false

	if errReadFile == nil {

		errUnmarshal := yaml.Unmarshal(yamlFile, &parsedConfig)

		if errUnmarshal != nil {

			log.Printf("WARN: Error parsing %s: %v. Proceeding with defaults.", configPath, errUnmarshal)

		} else {

			configExistsAndReadable = true

		}

	} else if !os.IsNotExist(errReadFile) { // File exists but other read error

		log.Printf("WARN: Error reading %s: %v. Proceeding with defaults.", configPath, errReadFile)

	}

	if configExistsAndReadable && parsedConfig.Auth.Username != "" && parsedConfig.Auth.Password != "" {

		AppCredentials.Username = parsedConfig.Auth.Username

		AppCredentials.Password = parsedConfig.Auth.Password

		log.Printf("INFO: Loaded credentials from %s for user: %s", configPath, AppCredentials.Username)

	} else {

		AppCredentials.Username = defaultUsername

		randomPassword, errGenPass := serverutil.GenerateRandomPassword(20)

		if errGenPass != nil {

			log.Fatalf("FATAL: Could not generate random password: %v", errGenPass)

		}

		AppCredentials.Password = randomPassword

		if os.IsNotExist(errReadFile) {

			log.Printf("INFO: Configuration file %s not found.", configPath)

		} else if !configExistsAndReadable {

			// This case is covered by earlier logs (read or unmarshal error)

		} else { // Config existed but was incomplete

			log.Printf("INFO: Incomplete or missing auth credentials in %s.", configPath)

		}

		log.Printf("INFO: Using default admin user. Username: %s, Generated Password: %s", AppCredentials.Username, AppCredentials.Password)

		log.Printf("INFO: To use custom credentials, set auth.username and auth.password in %s", configPath)

		// Write back the default/generated credentials to api.yaml

		currentConfigOnDisk := APIConfigYAML{}

		if configExistsAndReadable { // Preserve other data if file was read successfully

			_ = yaml.Unmarshal(yamlFile, &currentConfigOnDisk) // Re-unmarshal to get full existing structure if any

		}

		currentConfigOnDisk.Auth.Username = AppCredentials.Username

		currentConfigOnDisk.Auth.Password = AppCredentials.Password

		updatedYAML, marshalErr := yaml.Marshal(&currentConfigOnDisk)

		if marshalErr != nil {

			log.Printf("WARN: Could not marshal default config to YAML: %v", marshalErr)

		} else {

			writeErr := os.WriteFile(configPath, updatedYAML, 0644)

			if writeErr != nil {

				log.Printf("WARN: Could not write default config to %s: %v", configPath, writeErr)

			} else {

				log.Printf("INFO: Default/updated configuration written to %s", configPath)

			}

		}

	}

}
