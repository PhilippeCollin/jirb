package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func askSinglePreference(text, currentValue string) string {
	var pref string
	allowDefault := currentValue != ""
	if allowDefault {
		fmt.Printf("%s (%s): ", text, currentValue)
	} else {
		fmt.Printf("%s: ", text)
	}

	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	check(err)
	pref = string(line)

	if allowDefault && pref == "" {
		pref = currentValue
	}
	return pref
}

func askPreferences(prefs Preferences) Preferences {
	jiraURL := askSinglePreference("Jira Host URL", prefs.JiraHostURL)
	parsedURL, err := url.ParseRequestURI(jiraURL)
	check(err)
	prefs.JiraHostURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	return prefs
}

func updateAllPreferences() Preferences {
	prefs := readPrefsFileOrCreate()
	updatedPrefs := askPreferences(prefs)
	savePreferencesToFile(updatedPrefs)
	return updatedPrefs
}

func deletePreferencesFile() {
	err := os.Remove(getConfigFilePath())
	if !os.IsNotExist(err) {
		check(err)
	}
}

func savePreferencesToFile(prefs Preferences) {
	filePath := getConfigFilePath()
	newFileBytes, err := json.MarshalIndent(&prefs, "", "  ")
	check(err)
	err = ioutil.WriteFile(filePath, newFileBytes, 0660)
	check(err)
}

func readPrefsFileOrCreate() Preferences {
	filePath := getConfigFilePath()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0660)
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
		updatedPrefs := askPreferences(prefs)
		savePreferencesToFile(updatedPrefs)
		return updatedPrefs
	}

	return prefs
}

func getOrCreatePreferences() Preferences {
	var prefs Preferences
	if cachedPrefs != nil {
		prefs = *cachedPrefs
	} else {
		prefs = readPrefsFileOrCreate()
	}
	cachedPrefs = &prefs
	return prefs
}
