package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)
type processFile struct {
	name string
	path string
	count int
	port int
	nginx string
	comp string
}
var Reset   = "\033[0m"
var Red     = "\033[31m"
var Green   = "\033[32m"
func newTask(name string, path string, count int, port int, comp string, nginx string) *processFile {
	p:= processFile{name,path, count, port , nginx, comp}
	return &p
}
func readLine(fn string, ex string) (int, error) {
	var ret = -1
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line, dispose string
	line, err = bf.ReadString('\n')
	if err != nil {
		return ret, nil
	}
	for line != "" {
		if strings.Contains(line, ex) {
			fmt.Sscanf(line, "%s %d",&dispose, &ret)
			return ret, nil
		}
		line, err = bf.ReadString('\n')
	}
	return ret, err
}
func parseFiles (task *processFile, p chan <-processFile, wg *sync.WaitGroup) {
	defer wg.Done() //uncomment in wg != nil
	if port, err := readLine(task.path, "EXPOSE");
		err == nil {
		task.port = port
		fmt.Println(Green, "dockerfile: ", task.name , " ", port, Reset)
		p <- processFile{ task.name,task.path,task.count,port,task.nginx,task.comp}
		} else {
		fmt.Println(Red, task.name , Reset)
	}
	return
}

func main(){
	startTime := time.Now()
	pListFile := make([]*processFile, 0) //dockerfiles
	services := make([]*processFile, 0)	//services ->dockerfile:container_name
	f, err := os.Open("services.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	line, err := bf.ReadString('\n')
	if err != nil {
		return
	}
	//skip header
	line, err = bf.ReadString('\n')
	if err != nil {
		return
	}

	for line != "" {
		line = strings.ReplaceAll(line, "\n", "")
		str := strings.Split(line, "|")
		s, err := strconv.Atoi(str[0])
		if err != nil {
			fmt.Println("invalid syntax here : ", "\"", line, "\"")
			return
		}
		if len(str) == 5{
			pListFile = append(pListFile, newTask(str[1], str[2], s, -1,str[3], str[4]))
		} else {
			pListFile = append(pListFile, newTask(str[1], str[2], s,  -1,str[3], ""))
		}
		line, err = bf.ReadString('\n')
	}
	for _, service:= range pListFile {
		var count = 1
		if service.count == 1 {
			services = append(services, newTask(service.name,service.path, service.count, -1, service.comp, service.nginx))
		}else {
			for count <= service.count {
				name := service.name + strconv.Itoa(count)
				services = append(services, newTask(name,service.path, service.count,-1, service.comp, service.nginx))
				count++
			}
		}
	}
	fmt.Println(Green, "Services.txt read", Reset)
	/*
	services.txt read
		- parse Dockerfile and match port with name
		- if more than 1 instance, parse nginx and match name, count, port
		- parse docker-compose and match name,count,port,hostname,containername
	*/
	var wgOne sync.WaitGroup
	var wgTwo sync.WaitGroup
	var waitNum = 1
	serviceChan := make(chan processFile, len(services))	//services ->dockerfile:container_name
	for _, service:= range services {
		wgOne.Add(waitNum)
		//validate Dockerfiles
		go parseFiles(service, serviceChan, &wgOne)
	}
	wgOne.Wait()
	close(serviceChan)
	s := make([]*processFile, 0)
	i:=0
	for service := range serviceChan{
		s = append(s, newTask(service.name,service.path, service.count, service.port, service.comp, service.nginx))
		wgTwo.Add(waitNum)
		go parseCompose(s[i], &wgTwo)
		if service.nginx != "" {
			wgTwo.Add(waitNum)
			//validate respective nginx
			go parseNginx(s[i], &wgTwo)
		}
		i++
	}
	wgTwo.Wait()
	/*
		------------- SERVICES ARE ASSUMED TO BE RUNNING BELOW THIS LINE -------------

	*/
	//validate all services are running
	if len(os.Args) != 1 {
		validator(services)
	}
	fmt.Println("Time taken: ", time.Now().Sub(startTime).String())
	return
}