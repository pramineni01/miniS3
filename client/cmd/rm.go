/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		aliasName, _ := cmd.Flags().GetString("aliasName")
		bucketName, _ := cmd.Flags().GetString("bucketName")
		if (aliasName == "") || (bucketName == "") {
			log.Fatal("Invalid input.")
		}

		jsonBody, _ := json.Marshal(map[string]string{"bucketName": bucketName})
		req, err := getRequestWithCredentials(http.MethodDelete, "rm", aliasName, jsonBody)
		if err != nil {
			log.Fatal("Error constructing http request object: ", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf(string(resBody))

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.Flags().String("aliasName", "", "Alias name")
	rmCmd.Flags().String("bucketName", "", "Bucket name")
}
