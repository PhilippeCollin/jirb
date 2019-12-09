package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
)

const keyringKey = "jirbCreds"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getKeyring() keyring.Keyring {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "Jirb",
	})
	check(err)
	return ring
}

func askCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Jira username: ")
	username, err := reader.ReadString('\n')
	check(err)

	fmt.Println("Password: ")
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	check(err)
	password := string(passwordBytes)

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

	fmt.Println(encodedCredentials)

	request, err := http.NewRequest(http.MethodGet, "https://services.boreal-is.com/jira/api/rest/latest/issue/MOB-1", nil)
	check(err)

	client := &http.Client{}

	resp, err := client.Do(request)
	check(err)

	data map = nil;
	err = json.NewDecoder(resp.Body).Decode(&data)

	request.Header.Add("Authorization", "Basic "+encodedCredentials)
}
