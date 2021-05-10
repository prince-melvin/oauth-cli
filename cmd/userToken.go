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
var userTokenCmd = &cobra.Command{
	Use:   "usertoken",
	Short: "Get token using user credentials",
	Long: `Uses oauth's  resource owner password credentials flow to obtain the claims such as access token and id token
	Usage:
	 - oauth usertoken -u <username> -p <password>
	 - oauth usertoken -u <username> -p <password> --access-token
	 - oauth usertoken -u <username> -p <password> --a`,
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

		getToken(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(userTokenCmd)
	userTokenCmd.Flags().StringP("username", "u", "", "HSDP IAM Username")
	userTokenCmd.Flags().StringP("password", "p", "", "HSDP IAM Password")
	userTokenCmd.Flags().BoolP("access-token", "a", false, "get only access token")
	userTokenCmd.MarkFlagRequired("username")
	userTokenCmd.MarkFlagRequired("password")
}

func getToken(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	cfg, err := pkg.ReadConfigFile()
	if err != nil {
		fmt.Println("Error reading config file")
	}
	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("username", username)
	params.Add("password", password)

	result, err := requests.New(cfg.Provider.TokenUrl).
		WithMethod("POST").
		WithBody(bytes.NewBufferString(params.Encode())).
		SetHeader("Authorization", fmt.Sprintf("Basic %s", pkg.GetBasicAuthHeader())).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Do().UnmarshalJSON()

	if err != nil {
		fmt.Printf("Error while retrieving token %s", err)
		return
	}

	onlyAccessToken, _ := cmd.Flags().GetBool("access-token")
	if onlyAccessToken {
		token, _ := json.Marshal(result.Get("access_token"))
		fmt.Printf("Access Token - %s", string(token))
	} else {
		resBytes, _ := json.MarshalIndent(result, "", "\t")
		fmt.Println(string(resBytes))
	}

}
