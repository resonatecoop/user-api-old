package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-pg/migrations"

	"user-api/pkg/config"
	"user-api/pkg/postgres"
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

	cfgPath, err := filepath.Abs("./../../conf.local.yaml")
	if err != nil {
		exitf(err.Error())
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		exitf(err.Error())
	}

	psn := cfg.DB.Dev.PSN
	logQueries := cfg.DB.Dev.LogQueries
  timeoutSeconds := cfg.DB.Dev.TimeoutSeconds
	if testing {
		psn = cfg.DB.Test.PSN
		logQueries = cfg.DB.Test.LogQueries
		timeoutSeconds = cfg.DB.Test.TimeoutSeconds
	}
	db, err := pgsql.New(psn, logQueries, timeoutSeconds)
	if err != nil {
		exitf(err.Error())
	}

	oldVersion, newVersion, err := migrations.Run(db, flags...)
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
