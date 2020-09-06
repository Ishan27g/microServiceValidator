package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func verifyProcess(name string, target string) error{
	f, err := os.Open("dockerresult.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines,scanner.Text())
	}
	for _, line := range lines {
		if strings.Contains(line, name) { // change this to matches
			var n,p,q string
			fmt.Sscanf(line, "%s %s %s",&n, &p, &q)
			if (strings.Contains(p, target)|| strings.Contains(q, target)) && strings.Contains(line, "Up") && strings.Compare(n, name)==0 {
				fmt.Println(Green,  name ,Reset)
				return nil
			}else {
				errMsg := "image not valid : " + name + " " + target
				fmt.Println(Red,  name ,Reset)
				return errors.New(errMsg)
			}
		}
	}
	f.Close()
	return nil
}

func validateImage(imageName string, port int, wg *sync.WaitGroup)  {
	defer wg.Done()
	verifyProcess(imageName, strconv.Itoa(port))
}
