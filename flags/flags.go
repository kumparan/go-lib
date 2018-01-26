package flags

import (
	"flag"
	"log"
	"os"
)

type f interface {
	// Parse will parse flag object being set
	Parse(*flag.FlagSet, []string) error
}

func Parse(flg f) {
	err := flg.Parse(flag.CommandLine, os.Args[1:])
	if err != nil {
		// fatal because flag must successfully parsed
		log.Fatal("Failed to parse flag", err)
	}
}
