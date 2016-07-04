package main

import "flag"

func init() {
	flag.StringVar(&configInfo.ToSha, "ToSha", configInfo.ToSha, "Set commit hash up to which the change log should be generated")
}

func main() {
	flag.Parse()

	ci, err := apiConn.allCommits()

	if err != nil {
		panic(err)
	}

	sortedMsg, _ := sortedChangeLogMessages(ci)

	err = generateChangeLogHTML(sortedMsg, "ChangeLog.html")

	if err != nil {
		panic(err)
	}
}
