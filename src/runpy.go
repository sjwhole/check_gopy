package main

import (
	"os"
	"os/exec"
)

func RunPy() {
	//Windows
	current, _ := os.Getwd()
	cmd := exec.Command("C:\\Anaconda3\\envs\\check\\python.exe", current+"\\convert.py", current)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	_ = cmd.Wait()

	cmd = exec.Command("C:\\Anaconda3\\envs\\check\\python.exe", current+"\\excellib.py", current)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	_ = cmd.Wait()

}
