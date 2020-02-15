package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"syscall"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
)

const keyringKey = "jirbCreds"

func getKeyring() keyring.Keyring {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "Jirb",
	})
	check(err)
	return ring
}

func updateCredentials(ring keyring.Keyring) (string, string) {
	creds := retrieveExistingCredentials(ring)
	fmt.Printf("Jira username (%s): ", creds.Username)
	var newUsername string
	_, err := fmt.Scanln(&newUsername)
	check(err)
	if newUsername == "" {
		newUsername = creds.Username
	}
	var isDifferentUser = false
	if newUsername != creds.Username {
		isDifferentUser = true
	}

	var newPassword string
	fmt.Print("Password (Empty to make no change): ")
	_, err = fmt.Scanln(&newPassword)
	check(err)

	if isDifferentUser || newPassword != "" {
		saveCredentials(newUsername, newPassword, ring)
	}
	return newUsername, newPassword
}

func askCredentials() (string, string) {
	fmt.Print("Jira username: ")
	var username string
	_, err := fmt.Scanln(&username)
	check(err)

	fmt.Print("Password: ")
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(passwordBytes)
	check(err)

	return username, password
}

func encodeCredentials(username, password string) string {
	usernameWithPassword := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(usernameWithPassword))
}

func saveCredentials(username string, password string, ring keyring.Keyring) {
	data := KeychainData{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(data)
	check(err)
	err = ring.Set(keyring.Item{
		Key:         keyringKey,
		Data:        jsonData,
		Description: "Jira credentials used by the Jirb tool",
		Label:       "Jira Credentials",
	})
	check(err)
}

func retrieveExistingCredentials(ring keyring.Keyring) KeychainData {
	item, _err := ring.Get(keyringKey)
	check(_err)
	creds := KeychainData{}
	err := json.Unmarshal(item.Data, &creds)
	check(err)
	return creds
}
