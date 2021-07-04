package main

import (
	"flag"
)

const (
	diffFlag = "turns on the generation of the difference script (by default, an initialization script is generated)" +
		"accepts hash of commits  if only 1 commit was transferred," +
		"a commit is generated from the transferred one to the last"
	outFlag = "path to output file (by default init_script.sql)"
)

func parseFlags() (string, string) {
	diff := flag.String("diff", "false", diffFlag)
	out := flag.String("o", "", outFlag)

	flag.Parse()

	return *out, *diff
}
