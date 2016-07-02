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
        <title>Change log for {{.AppName}}</title>
    </head>
    <body>
        <h1>Change log for {{.AppName}} version {{.AppVersion}}</h1>
        {{- range $type, $scopeMessages := .Changes}}
            {{- if filter $type}}
            <h2>{{$type}}</h2>
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
		"perf":
		return true
	default:
		return false
	}
}

func generateChangeLogHTML(chg changes, filePath string) error {
	funcMap := template.FuncMap{
		"filter": filterTypes,
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
		Changes    changes
		AppName    string
		AppVersion string
	}{
		Changes:    chg,
		AppName:    configInfo.AppName,
		AppVersion: configInfo.AppVersion,
	})

	if err != nil {
		return err
	}

	return nil
}
