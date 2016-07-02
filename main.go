package main

func main() {
	ci, err := apiConn.allCommits()

	if err != nil {
		panic(err)
	}

	sortedMsg, _ := sortedChangeLogMessages(ci)

	//fmt.Printf("Sorted messages:\n%+v\n", sortedMsg)

	err = generateChangeLogHTML(sortedMsg, "ChangeLog.html")

	if err != nil {
		panic(err)
	}
}
