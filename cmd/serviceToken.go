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
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"oauth-cli/pkg"
	"oauth-cli/pkg/requests"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

const (
	BeginRsaKey      = "-----BEGIN RSA PRIVATE KEY-----"
	EndRsaKey        = "-----END RSA PRIVATE KEY-----"
	LineBreak        = "\n"
	GrantTypeWithJwt = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
	ErrNotRSAPublicKey     = errors.New("Key is not a valid RSA public key")
)

// servicetokencmd represents the userToken command
var serviceTokenCmd = &cobra.Command{
	Use:   "servicetoken",
	Short: "Get token using service account details",
	Long: `Converts service account details to JWT claims and exchanges it for access-token
	Usage:
	 - oauth serviceToken -s <serviceID> --private-key-file <path-to-privateKey>
	 - oauth serviceToken -s <serviceID> --private-key-file <path-to-privateKey> --jwt
	 oauth serviceToken -s <serviceID> --private-key-file <path-to-privateKey> --access-token`,
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
		cfg, err := pkg.ReadConfigFile()
		if err != nil {
			fmt.Println("Error reading config file")
		}
		getServiceToken(cmd, args, cfg)
	},
}

func init() {
	rootCmd.AddCommand(serviceTokenCmd)
	serviceTokenCmd.Flags().StringP("serviceId", "s", "", "Service Id")
	serviceTokenCmd.Flags().String("private-key-file", "", "Path to your private key")
	serviceTokenCmd.Flags().BoolP("jwt", "j", false, "get only JWT signed key")
	serviceTokenCmd.Flags().BoolP("access-token", "a", false, "get only JWT signed key")
	serviceTokenCmd.MarkFlagRequired("serviceId")
	serviceTokenCmd.MarkFlagRequired("serviceId")
}

func getServiceToken(cmd *cobra.Command, args []string, cfg *pkg.Configuration) {
	serviceId, _ := cmd.Flags().GetString("serviceId")
	privateKeyPath, _ := cmd.Flags().GetString("private-key-file")

	data, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Error reading private key, %v", err)
		return
	}
	privateKey := parsePrivateKey(data)
	jwt, err := CreateAndSignJWT(serviceId, privateKey, cfg)
	if err != nil {
		fmt.Printf("Error while signing jwt token %v", err)
		return
	}
	onlyJwt, _ := cmd.Flags().GetBool("jwt")
	if onlyJwt {
		fmt.Printf("Jwt signed key\n %s", jwt)
		return
	}

	result, err := ExchangeForToken(jwt, cfg)
	onlyAccessToken, _ := cmd.Flags().GetBool("access-token")
	if onlyAccessToken {
		token, _ := json.Marshal(result.Get("access_token"))
		fmt.Printf("Access Token - %s", string(token))
	} else {
		resBytes, _ := json.MarshalIndent(result, "", "\t")
		fmt.Println(string(resBytes))
	}

}

func CreateAndSignJWT(serviceId string, privateKey []byte, cfg *pkg.Configuration) (string, error) {
	var err error
	var block *pem.Block
	if block, _ = pem.Decode(privateKey); block == nil {
		return "", ErrKeyMustBePEMEncoded
	}
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return "", err
		}
	}
	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["aud"] = cfg.Provider.IssuerUrl
	claims["exp"] = now.Add(time.Hour * time.Duration(2)).Unix() // The expiration time after which the token must be disregarded.
	claims["iss"] = serviceId
	claims["sub"] = serviceId

	jwt, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(parsedKey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}
	return jwt, nil
}

func parsePrivateKey(privateKeyBytes []byte) []byte {
	privateKey := string(privateKeyBytes)
	if !strings.Contains(privateKey, "\n") &&
		strings.HasPrefix(privateKey, BeginRsaKey) &&
		strings.HasSuffix(privateKey, EndRsaKey) {
		privateKey = privateKey[:len(BeginRsaKey)] + LineBreak + privateKey[len(BeginRsaKey):]
		position := strings.Index(privateKey, EndRsaKey)
		privateKey = privateKey[:position] + "\n" + privateKey[position:]
	}
	return []byte(privateKey)
}

func ExchangeForToken(jwt string, cfg *pkg.Configuration) (*simplejson.Json, error) {

	params := url.Values{}
	params.Add("grant_type", GrantTypeWithJwt)
	params.Add("assertion", jwt)

	result, err := requests.New(cfg.Provider.TokenUrl).
		WithMethod("POST").
		WithBody(bytes.NewBufferString(params.Encode())).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Do().UnmarshalJSON()

	if err != nil {
		log.Fatalf("Error while exchanging token %s", err)
	}
	return result, nil
}
