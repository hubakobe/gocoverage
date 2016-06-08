package main

import (
	"flag"
	"gocoverage/coveragecaculator"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	caculator := coveragecaculator.NewCoveragecaculator(root)
	caculator.Caculate()
}
