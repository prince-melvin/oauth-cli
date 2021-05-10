package pkg

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	configFile = "providers.yaml"
)

type Configuration struct {
	Provider Provider
}
type Provider struct {
	Name          string
	ClientID      string
	ClientSecret  string
	TokenUrl      string
	IntrospectUrl string
	IssuerUrl     string
}

var HomeDirectory string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	HomeDirectory = home
}

func getFilePath() string {
	return path.Join(HomeDirectory, ".oauth", configFile)
}

func getConfigDirPath() string {
	return filepath.Dir(getFilePath())
}

//Check existence of the configuration file
func IsConfigFileExist() bool {
	if fi, err := os.Stat(getFilePath()); err != nil || fi.IsDir() {
		return false
	}
	return true
}

//Read the configuration
func ReadConfigFile() (*Configuration, error) {
	viper.SetConfigFile(getFilePath())
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg Configuration
	err = viper.UnmarshalKey("configuration", &cfg)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
	return &cfg, nil
}

//Write data of config to the configuration file
func WriteConfigFile(cfg *Configuration) error {
	createDirectoryIfNotExists()
	_, err := os.Create(getFilePath())
	if err != nil {
		return err
	}
	viper.SetConfigName("providers")
	viper.AddConfigPath(getConfigDirPath())
	viper.SetConfigType("yml")
	viper.Set("configuration", cfg)
	return viper.WriteConfig()
}

//Write default config into the configuration file
func WriteDefaultConfig() error {
	return WriteConfigFile(&Configuration{
		Provider: Provider{
			IntrospectUrl: "https://iam-client-test.us-east.philips-healthsuite.com/authorize/oauth2/introspect?api-version=3",
			TokenUrl:      "https://iam-client-test.us-east.philips-healthsuite.com/authorize/oauth2/token?api-version=2",
			IssuerUrl:     "https://iam-client-test.us-east.philips-healthsuite.com/oauth2/access_token",
			Name:          "HSDP",
			ClientID:      "<clientId>",
			ClientSecret:  "<clientsecret>",
		},
	})
}

func GetBasicAuthHeader() string {
	cfg, err := ReadConfigFile()
	if err != nil {
		fmt.Println("Error reading config file")
	}
	return b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.Provider.ClientID, cfg.Provider.ClientSecret)))
}

func createDirectoryIfNotExists() {
	if _, err := os.Stat(getConfigDirPath()); os.IsNotExist(err) {
		os.Mkdir(getConfigDirPath(), os.ModeDir)
	}
}
