package pkgcodelinecaculator

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Pkgcodelinecaculator struct {
	rootPath  string
	goFileTag string
}

func (pcl *Pkgcodelinecaculator) PackageFullPath() (fullPath string) {
	fullPath, _ = filepath.Abs(pcl.rootPath)
	return
}

func (pcl *Pkgcodelinecaculator) Caculate() (codeLine, commentLine int) {
	goFiles, _ := pcl.listFiles()
	for _, file := range goFiles {
		var cl int
		var coml int
		cl, coml = pcl.caculateFileCodeLine(file)
		codeLine += cl
		commentLine += coml
		fmt.Println(file)
	}
	return
}

func (pcl *Pkgcodelinecaculator) caculateFileCodeLine(fileName string) (codeLine, commentLine int) {
	f, err := os.Open(fileName)
	if nil != err {
		log.Println(err)
		return
	}
	defer f.Close()

	strComment := "//"

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), strComment) {
			commentLine += 1
			continue
		}
		codeLine += 1
	}
	return
}

func (pcl *Pkgcodelinecaculator) listFiles() ([]string, error) {
	var goFiles []string
	dir, err := ioutil.ReadDir(pcl.rootPath)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix := strings.ToUpper(pcl.goFileTag)

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		//		if sub := strings.Split(fi.Name(), "_test"); len(sub) > 1 {
		//			continue
		//		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			goFiles = append(goFiles, pcl.rootPath+PthSep+fi.Name())
		}
	}

	return goFiles, nil
}

func NewPkgcodelinecaculator(path string) *Pkgcodelinecaculator {
	return &Pkgcodelinecaculator{
		rootPath:  path,
		goFileTag: ".go",
	}
}
