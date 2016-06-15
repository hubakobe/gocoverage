package coveragecaculator

import (
	"flag"
	"fmt"
	"gocoverage/pkgcodelinecaculator"
	"gocoverage/pkgcoverageratecaculator"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type caculator struct {
	linecaculator *pkgcodelinecaculator.Pkgcodelinecaculator
	rater         *pkgcoverageratecaculator.Pkgcoverageratecaculator
}

type Coveragecaculator struct {
	rootPath   string
	caculators []caculator
}

func getAllPackage(path string) []string {
	var packageDirs []string
	var extractPaths []string
	args := flag.Args()
	args = args[1:]
	for _, arg := range args {
		if arg != "" {
			extractPaths = append(extractPaths, arg)
		}
	}

	filepath.Walk(path, func(packagePath string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		packageDirs = append(packageDirs, packagePath)
		return nil
	})
	for index, value := range packageDirs {
		if sub := strings.Split(value, ".git"); len(sub) > 1 {
			packageDirs[index] = ""
			continue
		}
		for _, path := range extractPaths {
			if sub := strings.Split(value, path); len(sub) > 1 {
				packageDirs[index] = ""
				fmt.Println(value, path)
				break
			}
		}
	}
	var packageDirsFilterd []string
	for _, value := range packageDirs {
		if value != "" {
			packageDirsFilterd = append(packageDirsFilterd, value)
		}
	}
	for _, value := range packageDirsFilterd {
		fmt.Println(value)
	}
	return packageDirsFilterd
}

func NewCoveragecaculator(rootPath string) *Coveragecaculator {
	packageDirs := getAllPackage(rootPath)
	var caculs []caculator
	for _, packageDir := range packageDirs {
		var c caculator
		c.linecaculator = pkgcodelinecaculator.NewPkgcodelinecaculator(packageDir)
		c.rater = pkgcoverageratecaculator.NewPkgcoverageratecaculator(packageDir)
		go c.rater.Caculate()
		caculs = append(caculs, c)
	}
	return &Coveragecaculator{
		rootPath:   rootPath,
		caculators: caculs,
	}
}

func (cc *Coveragecaculator) Caculate() {
	var totalLine int
	var totalPassLine int
	for _, caculator := range cc.caculators {
		var rate float64
		c1Line, comLine := caculator.linecaculator.Caculate()
		totalLine += c1Line
		rateString := <-caculator.rater.RateChannel()
		if rateString == "" {
			rate = 0
		} else {
			sub := strings.Split(rateString, ": ")
			if len(sub) < 2 {
				rate = 0 //usecase fail, then rate set 0
			} else {
				sub = strings.Split(sub[1], "%")
				rate, _ = strconv.ParseFloat(sub[0], 32)
			}
		}
		c1PassLine := c1Line * int(rate)
		totalPassLine += c1PassLine
		fmt.Println("===========================================================")
		fmt.Println("package path: ", caculator.linecaculator.PackageFullPath())
		fmt.Println("go file line: ", c1Line+comLine)
		fmt.Println("go code line: ", c1Line)
		fmt.Println("go comment line: ", comLine)
		fmt.Println("package coverage: ", rateString)
		fmt.Println("===========================================================")
	}
	fmt.Println("Total Code Line: ", totalLine)
	fmt.Printf("Total Coverage %.1f", float32(totalPassLine)/float32(totalLine))
	fmt.Println("%")

	return
}
