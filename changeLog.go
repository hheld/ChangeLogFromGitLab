package main

import (
	"html/template"
	"os"
)

type changeType string
type scope string
type scopeMessages map[scope][]string
type changes map[changeType]scopeMessages

const changeLogTemplate = `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Change log</title>
    </head>
    <body>
        <h1>Change log</h1>
        {{- range $type, $scopeMessages := .Changes}}
            {{- if filter $type}}
            <h2>{{$type | name}}</h2>
            <ul>
            {{- range $scope, $messages := $scopeMessages}}
                <li>
                    <h3>{{$scope}}</h3>
                    <ul>
                    {{- range $message := $messages}}
                        <li>{{$message}}</li>
                    {{- end}}
                    </ul>
                </li>
            {{- end}}
            </ul>
            {{- end}}
        {{- end}}
    </body>
</html>
`

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

func filterTypes(typeName changeType) bool {
	switch typeName {
	case "feat",
		"fix",
		"docs",
		"perf",
		"refactor":
		return true
	default:
		return false
	}
}

func typeName(typeName changeType) string {
	switch typeName {
	case "feat":
		return "New features"
	case "fix":
		return "Bug fixes"
	case "docs":
		return "Documentation"
	case "perf":
		return "Performance improvements"
	case "refactor":
		return "Internal changes"
	default:
		return string(typeName)
	}
}

func generateChangeLogHTML(chg changes, filePath string) error {
	funcMap := template.FuncMap{
		"filter": filterTypes,
		"name":   typeName,
	}

	tmpl, err := template.New("changeLog").Funcs(funcMap).Parse(changeLogTemplate)

	if err != nil {
		return err
	}

	f, err := os.Create(filePath)

	if err != nil {
		return err
	}

	err = tmpl.Execute(f, struct {
		Changes changes
	}{
		Changes: chg,
	})

	if err != nil {
		return err
	}

	return nil
}
