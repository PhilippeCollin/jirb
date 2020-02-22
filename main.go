package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
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
	updateConfig := flag.Bool("config", false, "Cycle through configurations and optionally change values.")
	reset := flag.Bool("reset", false, "Remove all configurations files and keychain entries, making it as if you had never run this tool.")
	help := flag.Bool("help", false, "Show this help message.")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	nFlags := flag.NFlag()
	if nFlags > 1 {
		panic("Cannot use more than one flag at the same time.")
	}

	if *reset {
		deletePreferencesFile()
		deleteCredentials()
	}

	if *updateConfig {
		updateAllPreferences()
		updateCredentials()
	}

	prefs := getOrCreatePreferences()
	creds := getOrCreateCredentials()
	basicAuthHeader := getBasicAuthHeader(creds.Username, creds.Password)

	request := prepareRequest(prefs.JiraHostURL, creds.Username)

	request.Header.Add("Authorization", basicAuthHeader)

	client := &http.Client{}
	resp, err := client.Do(request)
	if resp.StatusCode != 200 {
		fmt.Println("Received status %i from jira server", resp.Status)
	}
	check(err)
	defer resp.Body.Close()

	data := IssuesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	fmt.Println(resp.Status)
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

	fmt.Print("Branch name: ")
	input, err := readline.New("")
	check(err)
	input.Stdout().Write([]byte("Branch Name: "))
	input.WriteStdin([]byte(issueKind + "/" + branchName))
	branchName, err = input.Readline()
	check(err)

	cmd := exec.Command("git", "checkout", "-b", branchName)
	err = cmd.Run()
	check(err)
}
