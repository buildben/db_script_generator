package main

import (
	"flag"
)

const outFlag = "path to output file (by default init_script.sql)"

func parseFlags() string {
	out := flag.String("o", "", outFlag)
	flag.Parse()

	return *out
}
