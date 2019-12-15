package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

func mapIssuesToPromptItems(issues []Issue) []string {
	choices := make([]string, len(issues))
	for i, issue := range issues {
		choices[i] = fmt.Sprintf("%s - %s", issue.Key, issue.Fields.Summary)
	}
	return choices
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

	promptChoices := mapIssuesToPromptItems(data.Issues)
	prompt := promptui.Select{
		Label: "Select JIRA Issue",
		Items: promptChoices,
	}
	_, result, err := prompt.Run()
	check(err)
	branchName := formatBranchName(result)

	prompt = promptui.Select{
		Label: "Select Issue Kind",
		Items: []string{
			"Feature",
			"Hotfix",
		},
	}
	_, result, err = prompt.Run()
	check(err)
	issueKind := strings.ToLower(result)

	fmt.Println(branchName)

	fmt.Print("Branch name: ")
	input, err := readline.New("")
	check(err)
	input.Stdout().Write([]byte("Branch Name: "))
	input.WriteStdin([]byte(issueKind + "/" + branchName))
	branchName, err = input.Readline()
	check(err)

	fmt.Println(branchName)
}
