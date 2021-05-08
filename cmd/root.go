/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"oauth-cli/pkg"

	"github.com/spf13/cobra"
)

var (
	version    bool
	showConfig bool
)
var rootCmd = &cobra.Command{
	Use:   "oauth",
	Short: "Can be used to generate and inspect oauth2 tokens",
	Long: `Generic oauth2-compliant cli which can be used to generate user tokens using client credentials flow and 
	generate tokens based on your service account details`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !pkg.IsConfigFileExist() {
			err := pkg.WriteDefaultConfig()
			if err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			return printVersion()
		}
		if showConfig {
			return showCurrentConfig()
		}
		return cmd.Help()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "show current version of CLI")
	rootCmd.Flags().BoolVarP(&showConfig, "show config", "c", false, "show current configuration")
}

func printVersion() error {
	fmt.Println("Version: ", pkg.Version)
	return nil
}

func showCurrentConfig() error {
	cfg, err := pkg.ReadConfigFile()
	if err != nil {
		return err
	}
	fmt.Printf("Configured IDP: %s\n"+
		"Token-Url: %s\n"+
		"Introspect Url: %s\n"+
		"ClientId: %s", cfg.Provider.Name, cfg.Provider.TokenUrl, cfg.Provider.IntrospectUrl, cfg.Provider.ClientID)
	return nil
}
