package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func validateComp(fn string, ex string) int {
	var ret = 0
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err, fn)
		return -1
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	line, err := bf.ReadString('\n')
	if err != nil {
		return ret
	}
	for line != "" {
		line = strings.ReplaceAll(line, "\n", "")
		if strings.Contains(line, ex) {
			ret++
		}
		line, err = bf.ReadString('\n')
	}
	return ret
}
func parseCompose(task *processFile, wg *sync.WaitGroup)  {
	defer wg.Done()
	var count= validateComp(task.comp, task.name)
	if count == 3 || count == 4 {
		fmt.Println(Green, "docker-compose: ", task.name, Reset)
	}else {
		fmt.Println(Red, task.comp, task.name, Reset)
	}
}
