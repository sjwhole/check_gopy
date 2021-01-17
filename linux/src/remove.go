package main

import (
	"io/ioutil"
	"os"
	"regexp"
)

func RemoveFile(fileExtension string) {

	current, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	files, err10 := ioutil.ReadDir(current)
	if err10 != nil {
		panic(err10)
	}
	for _, f := range files {
		matched, _ := regexp.MatchString(fileExtension + "$", f.Name())
		if matched {
			_ = os.Remove(f.Name())
		}
	}
}
