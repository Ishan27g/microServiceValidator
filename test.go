package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

)
var Reset   = "\033[0m"
var Red     = "\033[31m"
var Green   = "\033[32m"
var results = "dockerresult.txt"
/*dockerCompose := [2]string{
	"../Service-Container/docker-compose.yaml",
	"./API-container/docker-compose.yaml"}


containers := []string{
	"api-nginx",
	"register_api",
	"post_api",
	"sync_api_1",
	"sync_api_2",
	"sync_api_3",
	"service-nginx",
	"register_service",
	"write_service",
	"sync_service1",
	"sync_service2",
	"sync_service3",
	"auth_service1",
	"auth_service2",
	"auth_service3",
	"db",
	"adminer"}
*/
var imageCount int = 17 + 1 // +1 for header
var images = [9]string{
	"register_service",
	"write_service",
	"sync_service",
	"service-nginx",
	"auth_service",
	"post_api",
	"register_api",
	"api-nginx",
	"sync_api",
}
/*var images = [9]string{
"service-register",
"service-write",
"service-sync",
"service-nginx",
"service-auth",
"api-post",
"api-register",
"api-nginx",
"api-sync",
}*/
var dockerfiles = [9]string{
"../Service-Container/service_register/Dockerfile",
"../Service-Container/service_write/Dockerfile",
"../Service-Container/service_sync/Dockerfile",
"../Service-Container/Service-container-nginx/Dockerfile",
"../Service-Container/service_auth/Dockerfile",
"../API-container/write/Dockerfile",
"../API-container/register/Dockerfile",
"../API-container/API-container-nginx/Dockerfile",
"../API-container/sync/Dockerfile"}
var targets = []string{
"EXPOSE"}
type process struct {
	name string
	path string
	port int
}
var pList [9]*process

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
func readExpose ( task *process, target string, wg *sync.WaitGroup)  {
	defer wg.Done() //uncomment in wg != nil
	if port, err := readLine(task.path, target);
		err == nil {
			task.port = port
		//	fmt.Println(task.name, "\t:\t", task.port)
	} else {
		fmt.Println("rsl:", err)
	}
}
func newProcess(name string, path string) *process {
	p:= process{name,path, -1}
	return &p
}
func checkProcessCount() (string) {
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

func verifyProcess(name string, target string) error{
	f, err := os.Open(results)
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
		if strings.Contains(line, name) {
			//fmt.Println("GOT : ", line)
			//fmt.Println("Expected value from dockerfiles : ", name, " : ",  target)
			if strings.Contains(line, target) && strings.Contains(line, "Up") {
				//fmt.Println("image valid", " : ", name, target)
				return nil
			}else {

				errMsg := "image not valid : " + name + " " + target
				return errors.New(errMsg)
			}
		}
	}
	f.Close()
	return nil
}

func getPortfromDockerfile(imagename string) int {
	for _, image := range pList{
		if strings.Compare(imagename, image.name) == 0{
			return image.port
		}
	}
	return -1
}
func logProcess() error {
	cmd := exec.Command("/bin/bash", "-c",
		"sudo docker ps  --format '{{.Names}}''->''{{.Ports}}''->''{{.Status}}' > dockerresult.txt")
	_, err := cmd.Output()
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	return nil
}
func validateCount() bool{
	count := strings.ReplaceAll(checkProcessCount(), "\n", "")
	if count != "" {
		if strings.Compare(count,"18") == 0 {
			return true
		} else{
			return false
		}
	}
	return false
}
func main() {

	num := 0
	var wgOne sync.WaitGroup
	var wgTwo sync.WaitGroup
	var waitNum = 1

	startTime := time.Now()
	//Fill processes structure
	for num<9 {
		pList[num]= newProcess(images[num], dockerfiles[num])
		num++
	}
	num = 0
	//For each file, launch a go-routine to extract port number from the corresponding file
	for num<9 {
		wgOne.Add(waitNum)
		go readExpose(pList[num], targets[0], &wgOne)
		num++
	}
	wgOne.Wait()
	fmt.Println("Docker files read")

	//Once all images are running, check process count
	if  true != validateCount() {

		fmt.Println(Red, "All services not running", Reset)
		oldTime := time.Now()
		diff := oldTime.Sub(startTime)
		fmt.Println("Time taken: ", diff.String())
		return
	}
	fmt.Println(Green, "All processes running",Reset)

	//For each file, launch a go-routine to verify if service is mapped to correct port
	logProcess()
	for _, image:= range images {
		wgTwo.Add(waitNum)
		go validateImage(image, &wgTwo)
	}
	wgTwo.Wait()
	oldTime := time.Now()
	diff := oldTime.Sub(startTime)
	fmt.Println("Time taken: ", diff.String())
	return
}
func validateImage(imageName string, wg *sync.WaitGroup)  {
	defer wg.Done()
	ports := getPortfromDockerfile(imageName)
	if ports != -1 {
		err := verifyProcess(imageName, strconv.Itoa(ports))
		if err != nil {
			fmt.Println("invalid port detected")
			fmt.Println(err)
		}
	}
}
