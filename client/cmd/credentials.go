package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const server_url = "http://localhost:8090"

type AliasInfo struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func saveCredsToDisk(aliasName string, accessKey string, secretKey string) error {
	// if creds file is already there, then append.
	// else add new file

	credsMap := loadCredsFromDisk(aliasName)
	fmt.Println(credsMap)
	credsMap[aliasName] = AliasInfo{AccessKey: accessKey, SecretKey: secretKey}

	f, err := os.OpenFile(fmt.Sprintf("./%s/json", aliasName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	byt, err := json.Marshal(credsMap)
	if err != nil {
		log.Fatal("Error marshaling creds map: ", err)
	}

	if _, err := f.Write(byt); err != nil {
		log.Fatal("Error writing to creds file: ", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return nil

}

func loadCredsFromDisk(aliasName string) map[string]AliasInfo {

	if _, err := os.Stat(fmt.Sprintf("./%s/json", aliasName)); err == nil {
		f, err := os.Open(fmt.Sprintf("./%s/json", aliasName))
		if err != nil {
			log.Fatal("Error opening creds file: ", err)
		}
		b, err := io.ReadAll(f)
		if err != nil {
			log.Fatal("Error reading creds file content: ", err)
		}
		defer f.Close()
		credsMap := map[string]AliasInfo{}
		if err := json.Unmarshal(b, &credsMap); err != nil {
			log.Fatal("Error unmarshaling creds file content.")
		}
		return credsMap
	}

	return map[string]AliasInfo{}
}

func getAccessSecretKeys(aliasName string) (*AliasInfo, error) {
	credsMap := loadCredsFromDisk(aliasName)
	ai, found := credsMap[aliasName]
	if !found {
		return nil, errors.New("Alias name not found.")
	}
	return &ai, nil
}

func getRequestWithCredentials(httpMethod string, endpoint string, aliasName string, body []byte) (*http.Request, error) {
	ai, err := getAccessSecretKeys(aliasName)
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/%s", server_url, endpoint)
	req, err := http.NewRequest(httpMethod, requestURL, bytes.NewReader(body))
	req.Header.Set("alias-name", aliasName)
	req.Header.Set("access-key", ai.AccessKey)
	req.Header.Set("secret-key", ai.SecretKey)

	return req, nil

}
