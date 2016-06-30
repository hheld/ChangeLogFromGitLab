package main

import "fmt"

func main() {
	ci, err := apiConn.allCommits()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Commits: %+v\nNumber of commits: %d\n", ci, len(ci))

	sortedMsg, _ := sortedChangeLogMessages(ci)

	fmt.Printf("Sorted messages:\n%+v\n", sortedMsg)
}
