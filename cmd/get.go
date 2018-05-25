// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"github.com/SDGophers/secret/vault"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires 1 args: key to get")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.NewVault(passphrase, "vault.encrypted")
		if err != nil {
			panic(fmt.Sprintf("failed to open vault: %v", err))
		}
		value, err := v.Get(args[0])
		if err != nil {
			panic(fmt.Sprintf("failed to retrieve value: %v", err))
		}
		fmt.Printf("value is: %s", value)
	},
}

var passphrase string

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.PersistentFlags().StringVarP(&passphrase, "key", "k", "", "Key used for opening vault")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
