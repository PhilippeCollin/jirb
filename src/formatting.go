package main

import (
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func formatBranchName(text string) string {
	r := regexp.MustCompile("([a-zA-Z]+-\\d+) - (.+)")

	matches := r.FindStringSubmatch(text)
	if len(matches) < 3 {
		panic("Unexpected ticket name format")
	}
	ticketID := matches[1]
	fmt.Println(matches[2])
	summary := matches[2]

	paranthesesRegexp := regexp.MustCompile("\\(.+\\)")
	summary = paranthesesRegexp.ReplaceAllString(summary, "")

	summary = strings.ToLower(summary)
	summary = strings.Join(strings.Fields(summary), " ")

	normalizedSummaryBytes := make([]byte, len(summary))
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, e := t.Transform(normalizedSummaryBytes, []byte(summary), true)

	summary = string(normalizedSummaryBytes)

	charsToRemove := regexp.MustCompile("[^a-zA-Z0-9 ]")
	summary = charsToRemove.ReplaceAllString(summary, "")

	summary = strings.Join(strings.Fields(summary), " ")
	summary = strings.ReplaceAll(summary, " ", "-")

	check(e)

	return ticketID + "_" + string(summary)
}
