package main

import (
	"fmt"
	"net/http"
)

func prepareRequest(jiraHostURL, username string) *http.Request {
	baseURL := fmt.Sprintf("%s/rest/api/latest/search", jiraHostURL)
	request, err := http.NewRequest(http.MethodGet, baseURL, nil)
	check(err)

	query := fmt.Sprintf("assignee=%s AND status=\"In Progress\"", username)
	q := request.URL.Query()
	q.Add("jql", query)

	request.URL.RawQuery = q.Encode()
	return request
}
