package main

// TODO OPT: Use a library/framework for CLI apps (check out Cobra)

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

const usage = `task is a CLI for managing your TODOs.

Usage:
  task [command]

Available Commands:
  add         Add a new task to your TODO list
  do          Mark a task on your TODO list as complete
  list        List all of your incomplete tasks
`

func add(desc string) {
}

func do(id int) {
}

func list() {
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println(usage)
		return
	}

	switch args[0] {
	case "add":
		add(strings.Join(args[1:], ""))
	case "do":
		if len(args) != 2 {
			fmt.Println(usage)
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		do(id)
	case "list":
		if len(args) > 1 {
			fmt.Println(usage)
			return
		}
		list()
	default:
		fmt.Println(usage)
	}
}
