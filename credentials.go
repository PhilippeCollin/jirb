package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
)

const keyringKey = "jirbCreds"

var cachedRing keyring.Keyring = nil

func getKeyring() keyring.Keyring {
	if cachedRing != nil {
		return cachedRing
	}
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "Jirb",
	})
	check(err)
	cachedRing = ring
	return ring
}

func updateCredentials() (string, string) {
	creds := retrieveExistingCredentials()

	fmt.Printf("Jira username (%s): ", creds.Username)
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	check(err)
	newUsername := string(line)

	if newUsername == "" {
		newUsername = creds.Username
	}
	var isDifferentUser = false
	if newUsername != creds.Username {
		isDifferentUser = true
	}

	fmt.Print("Password (Empty to make no change): ")
	line, _, err = reader.ReadLine()
	check(err)
	newPassword := string(line)

	if isDifferentUser || newPassword != "" {
		saveCredentialsToKeyring(newUsername, newPassword)
	}
	return newUsername, newPassword
}

func deleteCredentials() {
	ring := getKeyring()
	err := ring.Remove(keyringKey)
	if err != keyring.ErrKeyNotFound {
		check(err)
	}
}

func askCredentials() (string, string) {
	fmt.Print("Jira username: ")
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	check(err)
	username := string(line)

	fmt.Print("Password: ")
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(passwordBytes)
	check(err)

	return username, password
}

func getBasicAuthHeader(username, password string) string {
	usernameWithPassword := fmt.Sprintf("%s:%s", username, password)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(usernameWithPassword))
}

func saveCredentialsToKeyring(username string, password string) {
	ring := getKeyring()
	data := Credentials{
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

func retrieveExistingCredentials() Credentials {
	ring := getKeyring()
	item, err := ring.Get(keyringKey)
	creds := Credentials{}
	if err == nil {
		err = json.Unmarshal(item.Data, &creds)
		check(err)
	} else if err != keyring.ErrKeyNotFound {
		check(err)
	}
	return creds
}

func getOrCreateCredentials() Credentials {
	existingCreds := retrieveExistingCredentials()
	if !existingCreds.isComplete() {
		username, password := askCredentials()
		existingCreds.Username = username
		existingCreds.Password = password
		saveCredentialsToKeyring(username, password)
	}
	return existingCreds
}
