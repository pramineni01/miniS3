/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var aliasName, accessKey, secretKey, endpoint string

		aliasName, _ = cmd.Flags().GetString("aliasName")
		endpoint, _ = cmd.Flags().GetString("url")
		accessKey, _ = cmd.Flags().GetString("accessKey")
		secretKey, _ = cmd.Flags().GetString("secretKey")

		if (aliasName == "") || (accessKey == "") || (secretKey == "") || (endpoint == "") {
			log.Fatal("Missing required arguments.")
		}
		data, err := json.Marshal(map[string]string{
			"aliasName": aliasName,
			"accessKey": accessKey,
			"secretKey": secretKey,
		})
		if err != nil {
			log.Fatal("json marshal error")
		}
		r, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
		defer r.Body.Close()
		io.Copy(os.Stdout, r.Body)

		// write to local client storage
		saveCredsToDisk(aliasName, accessKey, secretKey)
	},
}

func init() {
	aliasCmd.AddCommand(setCmd)

	setCmd.Flags().String("aliasName", "", "Alias name")
	setCmd.Flags().String("accessKey", "", "Access key")
	setCmd.Flags().String("secretKey", "", "Secret key")
	setCmd.Flags().String("url", "", "MiniS3 http server endpoint")
}
