package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func logProcess() error {
	cmd := exec.Command("/bin/bash", "-c",
		"sudo docker ps  --format '{{.Names}}'' ''{{.Ports}}'' ''{{.Status}}' > dockerresult.txt")
	_, err := cmd.Output()
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	return nil
}
func processCount() (string) {
	ret := ""
	cmd := exec.Command("/bin/bash", "-c", "sudo docker ps | wc -l")
	op, err := cmd.Output()
	if err != nil {
		fmt.Print(err.Error())
		return ret
	}
	op = bytes.Trim(op, " ")
	ou := string(op)
	if ou != "" {
		return ou
	}
	return ret
}
func validateCount() bool{
	count := strings.ReplaceAll(processCount(), "\n", "")
	if count != "" {
		if strings.Compare(count,"18") == 0 {
			return true
		} else{
			return false
		}
	}
	return false
}
