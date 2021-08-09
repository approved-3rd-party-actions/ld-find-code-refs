package element

import "strings"

type ElementMatcher struct {
	Elements   []string
	Aliases    map[string][]string
	Delimiters []string
	ProjKey    string
	Directory  string
	Built      map[string][]string
}

type Matcher struct {
	Elements   []ElementMatcher
	Type       string
	CtxLines   int
	Delimiters string
}

func (m Matcher) MatchElement(line, flagKey string) bool {
	if m.Delimiters == "" && strings.Contains(line, flagKey) {
		return true
	}

	firstMatcher := m.Elements[0]
	delimitedFlag := firstMatcher.Built[flagKey]

	for _, flagKey := range delimitedFlag {
		if strings.Contains(line, flagKey) {
			return true
		}
	}
	return false
}
