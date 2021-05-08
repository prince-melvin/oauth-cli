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
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"oauth-cli/pkg"
	"oauth-cli/pkg/requests"

	"github.com/spf13/cobra"
)

// userTokenCmd represents the userToken command
var introspectionCmd = &cobra.Command{
	Use:   "introspect",
	Short: "Get introspection details for a token",
	Long: `Uses Iam's introspect api to fetch the access-token's claims and permssions
	Usage:
	 - oauth introspect <accesstoken>`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !pkg.IsConfigFileExist() {
			err := pkg.WriteDefaultConfig()
			if err != nil {
				return err
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		introspect(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(introspectionCmd)
}

func introspect(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		fmt.Println("Incorrect Usage, pass acces token")
		cmd.Help()
		return
	}

	cfg, err := pkg.ReadConfigFile()
	if err != nil {
		fmt.Println("Error reading config file")
	}

	token := args[0]
	params := url.Values{}
	params.Add("token", token)
	result, err := requests.New(cfg.Provider.IntrospectUrl).
		WithMethod("POST").
		WithBody(bytes.NewBufferString(params.Encode())).
		SetHeader("Authorization", fmt.Sprintf("Basic %s", pkg.GetBasicAuthHeader())).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Do().UnmarshalJSON()

	if err != nil {
		fmt.Printf("Error while intropecting token %s", err)
	}
	resBytes, _ := json.MarshalIndent(result, "", "\t")
	fmt.Println(string(resBytes))

}
