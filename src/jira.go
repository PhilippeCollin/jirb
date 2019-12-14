package main

import (
	"fmt"
	"net/http"
)

func formatRequestURL(jiraHostURL, username string) string {
	return fmt.Sprintf("%s/rest/api/latest/search?jql=assignee=%s AND status=\"In Progress\" AND status changed TO \"In Progress\" after -1w", jiraHostURL, username)
}

func prepareRequest(jiraHostURL, username string) *http.Request {
	baseURL := fmt.Sprintf("%s/rest/api/latest/search", jiraHostURL)
	request, err := http.NewRequest(http.MethodGet, baseURL, nil)
	check(err)

	query := fmt.Sprintf("assignee=%s AND status=\"In Progress\" AND status changed TO \"In Progress\" after -1w", username)
	q := request.URL.Query()
	q.Add("jql", query)

	request.URL.RawQuery = q.Encode()
	return request
}
