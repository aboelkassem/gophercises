package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	filenamePattern       *regexp.Regexp
	renamedFilenameFormat string
)

func main() {
	var (
		flagDir     = flag.String("dir", ".", "the path to the directory to work with")
		flagPattern = flag.String("pattern", "(.*?)\\s*\\((\\d+).+\\)", "the regx pattern to match filenames")
		flagFormat  = flag.String("format", "Episode %03s - %s", "the format to use for renaming")
	)
	flag.Parse()
	fmt.Println(*flagPattern)
	filenamePattern = regexp.MustCompile(*flagPattern)
	renamedFilenameFormat = *flagFormat

	if err := filepath.Walk(*flagDir, walkFn); err != nil {
		log.Fatal(err)
	}
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// return filepath.SkipDir // to skip this dir or file

	if info.IsDir() {
		return nil
	}

	fileExtension := filepath.Ext(info.Name())
	filename := strings.TrimSuffix(info.Name(), fileExtension)
	stringMatches := filenamePattern.FindStringSubmatch(filename)
	if len(stringMatches) == 0 {
		return nil
	}

	// convert []strings to []interface{}
	matches := make([]interface{}, len(stringMatches)-1)
	for i := 1; i < len(stringMatches); i++ {
		matches[i-1] = stringMatches[i]
	}
	renamedFilename := fmt.Sprintf(renamedFilenameFormat, matches...) + fileExtension
	renamedPath := filepath.Join(filepath.Dir(path), renamedFilename)

	fmt.Printf("old=%q, new=%q\n",
		path,
		renamedPath,
	)

	return os.Rename(path, renamedPath)
}
