/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	var aliasName, accessKey, secretKey, endpoint string

	// 	aliasName, _ = cmd.Flags().GetString("aliasName")
	// 	endpoint, _ = cmd.Flags().GetString("url")
	// 	accessKey, _ = cmd.Flags().GetString("accessKey")
	// 	secretKey, _ = cmd.Flags().GetString("secretKey")

	// 	if (aliasName == "") || (accessKey == "") || (secretKey == "") || (endpoint == "") {
	// 		log.Fatal("Missing required arguments.")
	// 	}
	// 	data, err := json.Marshal(map[string]string{
	// 		"aliasName": aliasName,
	// 		"accessKey": accessKey,
	// 		"secretKey": secretKey,
	// 	})
	// 	if err != nil {
	// 		log.Fatal("json marshal error")
	// 	}
	// 	r, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	// 	defer r.Body.Close()
	// 	io.Copy(os.Stdout, r.Body)
	// },
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
