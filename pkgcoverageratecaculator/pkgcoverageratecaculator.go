package pkgcoverageratecaculator

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type Pkgcoverageratecaculator struct {
	packagePath string
	goBuildPath string
	rateResult  chan string
}

func GetSrcFullPath(path string) (fullPath string) {
	fullPath, _ = filepath.Abs(path)
	return
}

func getGoBuildPath(path string) string {
	subString := strings.Split(path, "src/")
	if len(subString) == 1 {
		return subString[0]
	}
	return subString[1]
}

func (cr *Pkgcoverageratecaculator) RateChannel() chan string {
	return cr.rateResult
}

func (cr *Pkgcoverageratecaculator) hasTestFile() bool {
	dir, err := ioutil.ReadDir(cr.packagePath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	suffix := "_test.go"
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(fi.Name(), suffix) {
			return true
		}
	}

	return false
}
func (cr *Pkgcoverageratecaculator) Caculate() {
	if !cr.hasTestFile() {
		cr.rateResult <- ""
		return
	}
	result, err := exec.Command("go", "test", "-cover", cr.goBuildPath).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	strResult := string(result)
	cr.rateResult <- strResult
}

func NewPkgcoverageratecaculator(packagePath string) *Pkgcoverageratecaculator {
	builePath := getGoBuildPath(packagePath)
	fmt.Println("------" + builePath)
	return &Pkgcoverageratecaculator{
		packagePath: packagePath,
		goBuildPath: builePath,
		rateResult:  make(chan string),
	}
}
