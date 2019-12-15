package main

import (
	"encoding/base64"
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

func saveCredentials(encodedCredentials string, ring keyring.Keyring) {
	err := ring.Set(keyring.Item{
		Key:         keyringKey,
		Data:        []byte(encodedCredentials),
		Description: "Jira credentials used by the Jirb tool",
		Label:       "Jira Credentials",
	})
	check(err)
}
