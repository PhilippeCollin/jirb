package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

const configFileName = ".jirb"

var cachedPrefs *Preferences = nil

func getConfigFilePath() string {
	homedir, err := os.UserHomeDir()
	check(err)
	return fmt.Sprintf("%s/%s", homedir, configFileName)
}

func askSinglePreference(text string) string {
	var pref string
	fmt.Print(text)
	_, err := fmt.Scanln(&pref)
	check(err)
	return pref
}

func askMissingPreferences(prefs *Preferences) {
	if prefs.JiraHostURL == "" {
		jiraURL := askSinglePreference("Jira Host URL: ")
		parsedURL, err := url.ParseRequestURI(jiraURL)
		check(err)
		prefs.JiraHostURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	}

	if prefs.Username == "" {
		prefs.Username = askSinglePreference("Jira username: ")
	}
}

func readPrefsFileOrCreate() Preferences {
	filePath := getConfigFilePath()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0660)
	check(err)
	defer file.Close()

	stat, err := file.Stat()
	check(err)

	prefs := Preferences{}
	if stat.Size() > 0 {
		err := json.NewDecoder(file).Decode(&prefs)
		check(err)
	}

	if !prefs.isComplete() {
		askMissingPreferences(&prefs)
		newFileBytes, err := json.MarshalIndent(&prefs, "", "  ")
		check(err)
		_, err = file.Write(newFileBytes)
		check(err)
		defer file.Sync()
	}

	return prefs
}

func getPreferences() Preferences {
	var prefs Preferences
	if cachedPrefs != nil {
		prefs = *cachedPrefs
	} else {
		prefs = readPrefsFileOrCreate()
	}
	cachedPrefs = &prefs
	return prefs
}
