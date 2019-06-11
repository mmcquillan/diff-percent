package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kylelemons/godebug/diff"
	"github.com/mmcquillan/matcher"
)

func main() {

	// vars
	var dir1 string
	var dir2 string
	var filter string
	var exclude string
	var fileCount int
	var totalCount int
	var diffCount int

	// input
	match, _, values := matcher.Matcher("<bin> <path1> <path2> [--filter] [--exclude]", strings.Join(os.Args, " "))
	if match {
		dir1 = values["path1"]
		dir2 = values["path2"]
		if values["filter"] != "" && values["filter"] != "true" {
			filter = values["filter"]
		}
		if values["exclude"] != "" && values["exclude"] != "true" {
			exclude = values["exclude"]
		}
	} else {
		fmt.Println("diff-percent <path1> <path2> [--filter] [--exclude]")
		os.Exit(1)
	}

	// starting up
	fmt.Println("Diff Percent")

	// initialize
	file1 := GetFiles(dir1, filter, exclude)
	file2 := GetFiles(dir2, filter, exclude)
	fileCount = 0
	totalCount = 0
	diffCount = 0

	// compare from 1 to 2
	for k1, v1 := range file1 {
		fileCount++
		lc := LineCount(v1)
		totalCount += lc
		if v2, ok := file2[k1]; ok {
			d := Diff(v1, v2)
			if d > 0 {
				diffCount += d
				fmt.Print("[x]")
			} else {
				fmt.Print("[=]")
			}
		} else {
			diffCount += lc
			fmt.Print("[+]")
		}
		fmt.Println(" " + k1)
	}

	// compare from 2 to 1
	for k2, v2 := range file2 {
		if _, ok := file1[k2]; !ok {
			fileCount++
			lc := LineCount(v2)
			totalCount += lc
			diffCount += lc
			fmt.Println("[-] " + k2)
		}
	}

	// final output
	pd := (float64(diffCount) / float64(totalCount)) * 100
	fmt.Println("")
	fmt.Printf("      FILE COUNT: %d\n", fileCount)
	fmt.Printf(" LINE DIFF COUNT: %d\n", diffCount)
	fmt.Printf("LINE TOTAL COUNT: %d\n", totalCount)
	fmt.Printf("    PERCENT DIFF: %f\n", pd)

}

func Diff(f1 string, f2 string) (count int) {

	file1, err := ioutil.ReadFile(f1)
	if err != nil {
		log.Fatal(err)
	}

	file2, err := ioutil.ReadFile(f2)
	if err != nil {
		log.Fatal(err)
	}

	d := diff.Diff(string(file1), string(file2))

	scanner := bufio.NewScanner(strings.NewReader(d))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "+") {
			count++
		}
		if strings.HasPrefix(line, "-") {
			count++
		}
	}

	return count
}

func LineCount(filePath string) (count int) {
	file, _ := os.Open(filePath)
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	count = 0
	for fileScanner.Scan() {
		count++
	}
	return count
}

func GetFiles(dir string, filter string, exclude string) (files map[string]string) {
	files = make(map[string]string)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.Contains(f.Name(), filter) && (!strings.Contains(path, exclude) || exclude == "") {
			files[strings.Replace(path, dir, "", 1)] = path
		}
		return nil
	})
	if err != nil {
		log.Fatal("Walking path for dir " + dir + " : " + err.Error())
	}
	return files
}
