package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jgsheppa/localPass/generator"
)

func main() {
	passwordCmd := flag.NewFlagSet("password", flag.ExitOnError)
	length := passwordCmd.Int("length", 24, "Length of password")

	passwordCmd.Usage = func() {
		fmt.Println("Usage: [password] [...flags] ")
		fmt.Println("Flags:")
		passwordCmd.PrintDefaults()
	}

	if len(os.Args) < 2 {
		passwordCmd.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "password":
		passwordCmd.Parse(os.Args[2:])
		code, err := generator.Run(length)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(code)
	default:
		passwordCmd.Usage()
		os.Exit(1)
	}

}
