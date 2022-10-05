package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const MAX = 3
const TIMEOUT = 5

func main() {
	showIntro()
	for {
		showMenu()
		command := getCommand()
		fmt.Println("")
		switch command {
		case 1:
			startMonitor()
		case 2:
			showLogs()
		case 0:
			fmt.Println("Closed program...")
			os.Exit(0)
		default:
			fmt.Println("Command is not found")
			os.Exit(-1)
		}
	}
}

func showIntro() {
	version := 1.1
	fmt.Println("Welcome to Website Monitoring")
	fmt.Println("This program is in version:", version)
}

func showMenu() {
	fmt.Println("1 - Start Monitoring")
	fmt.Println("2 - View Logs")
	fmt.Println("0 - Closed Program")
}

func getCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func startMonitor() {
	fmt.Println("Monitoring...")
	urls := readUrlsFile()
	for count := 0; count < MAX; count++ {
		for index, url := range urls {
			fmt.Println("Testing the url", index, ":", url)
			testUrl(url)
		}
		fmt.Println("")
		time.Sleep(TIMEOUT * time.Second)
	}
	fmt.Println("")
}

func readUrlsFile() []string {
	var urls []string
	file, err := os.Open("urls.txt")
	if err != nil {
		fmt.Println(err)
	}
	read := bufio.NewReader(file)
	for {
		line, err := read.ReadString('\n')
		line = strings.TrimSpace(line)
		urls = append(urls, line)
		if err == io.EOF {
			break
		}
	}
	file.Close()
	return urls
}

func testUrl(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode == 200 {
		fmt.Println("Url", url, ": online")
		registerLogs(url, true)
	} else {
		fmt.Println("Url", url, "is with problem, status:", response.StatusCode)
		registerLogs(url, false)
	}
}

func registerLogs(url string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	timeCurrent := time.Now().Format("02/01/2006 15:04:05")
	file.WriteString(timeCurrent + " - " + url + " - online: " + strconv.FormatBool(status) + "\n")
	file.Close()
}

func showLogs() {
	fmt.Println("Show logs...")
	file, err := os.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(file))
}
