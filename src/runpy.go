package main

import (
	"os"
	"os/exec"
)

func main() {
	//Windows
	current, _ := os.Getwd()
	cmd := exec.Command("C:\\Anaconda3\\envs\\check\\python.exe", current+"\\convert.py", current)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	_ = cmd.Wait()

	//Linux
	//Use gnumeric ssconvert

	//current, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//files, err10 := ioutil.ReadDir(current)
	//if err10 != nil {
	//	panic(err10)
	//}
	//for _, f := range files {
	//	matched, _ := regexp.MatchString("xls$", f.Name())
	//	if matched {
	//		cmd := exec.Command("ssconvert", current, current)
	//		err := cmd.Start()
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		_ = cmd.Wait()
	//	}
	//}

	cmd = exec.Command("C:\\Anaconda3\\envs\\check\\python.exe", current+"\\excellib.py", current)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	_ = cmd.Wait()



}
