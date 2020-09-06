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

func validateNginx(fn string, ex string, port int) (int, error) {
	var ret = -1
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line, dispose , name string
	line, err = bf.ReadString('\n')
	if err != nil {
		return ret, nil
	}
	for line != "" {
		if strings.Contains(line, ex) {
			line = strings.ReplaceAll(line, ";", "")
			fmt.Sscanf(line, "%s %s",&dispose, &name)
			dispose = ex + ":" + strconv.Itoa(port)
			if strings.Compare(dispose, name) != 0 {
				return 0, errors.New("invalid: " + fn + " " + ex)
			}
			return ret, nil
		}
		line, err = bf.ReadString('\n')
	}
	return ret, err
}
func parseNginx(task *processFile, wg *sync.WaitGroup)  {
	defer wg.Done()
	if port, err := validateNginx(task.nginx, task.name, task.port);
		err == nil {
		task.port = port
		fmt.Println(Green, "nginx: ", task.name , Reset)
	}else {
		fmt.Println(Red, "nginx: ", task.name , Reset)
	}
}
