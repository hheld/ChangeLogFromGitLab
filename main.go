package main

//import "fmt"

func main() {
	//ci, err := apiConn.allCommits()
    _, err := apiConn.allCommits()

	if err != nil {
		panic(err)
	}

	//fmt.Printf("Commits: %+v\nNumber of commits: %d\n", ci, len(ci))
}
