package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/migrations"

  "user-api/internal/database"
)

const usageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
  go run *.go <command> [args]
`

func main() {
	flag.Usage = usage
	flag.Parse()

	var testing bool
	flags := flag.Args()
	if contains(flag.Args(), "testing") {
		testing = true
		flags = flags[:len(flags)-1]
	}
	fmt.Println(flags)

	oldVersion, newVersion, err := migrations.Run(database.Connect(testing), flags...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Printf(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
		   return true
		}
	}
	return false
}
