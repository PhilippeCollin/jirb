package main

type Preferences struct {
	JiraHostURL string
}

func (p Preferences) isComplete() bool {
	return p.JiraHostURL != ""
}

type Issue struct {
	Id     string
	Key    string
	Fields struct {
		Summary string
	}
}

type IssuesResponse struct {
	Total      int
	MaxResults int
	Issues     []Issue
}

type Credentials struct {
	Username string
	Password string
}

func (k Credentials) isComplete() bool {
	return k.Username != "" && k.Password != ""
}
