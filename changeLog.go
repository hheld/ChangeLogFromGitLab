package main

type changeType string
type scope string
type scopeMessages map[scope][]string
type changes map[changeType]scopeMessages

func isStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

func sortedChangeLogMessages(commits []commit) (changes, error) {
	chg := make(changes)

	for _, commit := range commits {
		desc := commit.ChageDesc

		if _, typeExists := chg[changeType(desc.Type)]; !typeExists {
			chg[changeType(desc.Type)] = make(scopeMessages)
		}

		scopeMap := chg[changeType(desc.Type)]

		if !isStringInSlice(desc.Subject, scopeMap[scope(desc.Scope)]) {
			scopeMap[scope(desc.Scope)] = append(scopeMap[scope(desc.Scope)], desc.Subject)
		}
	}

	return chg, nil
}
