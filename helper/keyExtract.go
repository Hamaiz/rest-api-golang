package helper

import (
	"regexp"
	"strings"
)

// KeyExtract - extract keywords from string
func KeyExtract(t string) string {
	// trim the string value
	t = strings.TrimSpace(t)

	// make string lowercase
	t = strings.ToLower(t)

	// remove charcters from string
	reg, _ := regexp.Compile(`[^a-zA-Z0-9\s]+]`)
	s := reg.ReplaceAllString(t, "")

	// convert string to array
	words := strings.Fields(s)

	// check if length is 0
	if len(words) == 0 {
		return ""
	}

	// remove all the duplicates
	words = unique(words)

	// remove all the stopwords
	words = check(words)

	// join array with pipe(|)
	t = strings.Join(words, "|")

	return t
}

// unique - gives only unique items in array
func unique(words []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	//checking for duplicates via loop
	for _, entry := range words {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

// check - checks and removes stopwords from array
func check(words []string) []string {
	// put all the stopwords into the map
	keys := make(map[string]bool)
	for _, stopword := range StopWords {
		keys[stopword] = true
	}

	// holds all the useful words
	list := []string{}

	// loop through all the words to check for useful words
	for _, entry := range words {
		if _, value := keys[entry]; !value {
			list = append(list, entry)
		}
	}

	return list
}
