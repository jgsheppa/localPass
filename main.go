package main

import (
	"fmt"
	"os"

	"github.com/jgsheppa/localPass/generator"
	"github.com/jgsheppa/localPass/list"
	"github.com/jgsheppa/localPass/pass"
)

func main() {
	passwordCmd := generator.Flag()

	if len(os.Args) < 2 {
		passwordCmd.FlagSet.Usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "password":
		passwordCmd.FlagSet.Parse(os.Args[2:])
		code, err := generator.Run(passwordCmd.Length)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(code)
	case "create":
		code, err := pass.Run()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(code)
	case "list":
		code, err := list.Run()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(code)
	default:
		passwordCmd.FlagSet.Usage()
		os.Exit(0)
	}
}
