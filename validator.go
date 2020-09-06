package main

import (
	"fmt"
	"sync"
)

func validator (services []*processFile) {
	var wgTwo sync.WaitGroup
	var waitNum = 1
	if  true != validateCount(){
		fmt.Println(Red, "Not all services are running", Reset)
		//return
	}
	logProcess()
	for _, service:= range services {
		wgTwo.Add(waitNum)
		go validateImage(service.name, service.port, &wgTwo)
	}
	wgTwo.Wait()
}
