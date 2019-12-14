package main

type Preferences struct {
	JiraHostURL string
	Username    string
}

func (p Preferences) isComplete() bool {
	return p.JiraHostURL != "" && p.Username != ""
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
