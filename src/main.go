package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/keyring"
)

const keyringKey = "jirbCreds"

func getKeyring() keyring.Keyring {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "Jirb",
	})
	check(err)
	return ring
}

func main() {
	fmt.Println("Hello, world.")

	ring := getKeyring()

	credsItem, err := ring.Get(keyringKey)
	encodedCredentials := ""
	if os.IsNotExist(err) || credsItem.Data == nil {
		username, password := askCredentials()
		encodedCredentials = encodeCredentials(username, password)
		saveCredentials(encodedCredentials, ring)
	} else {
		check(err)
		encodedCredentials = string(credsItem.Data)
	}

	preferences := getPreferences()
	request := prepareRequest(preferences.JiraHostURL, preferences.Username)
	check(err)

	request.Header.Add("Authorization", "Basic "+encodedCredentials)

	client := &http.Client{}
	resp, err := client.Do(request)
	check(err)
	defer resp.Body.Close()

	data := IssuesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	check(err)

	fmt.Println(data)
}
