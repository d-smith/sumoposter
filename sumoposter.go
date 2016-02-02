package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup
var messageCount int

func main() {
	//Grab the endpoint address
	endpoint := os.Getenv("SUMOENDPOINT")
	if endpoint == "" {
		fmt.Fprintf(os.Stderr, "Error - SUMOENDPOINT environment variable not set\n")
		os.Exit(1)
	}

	//Grab the sumo name, if any
	sumoName := os.Getenv("SUMONAME")

	//Check the args
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[1])
		os.Exit(1)
	}

	//Open the file pass to us by fluentd
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Read the file and process each line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go postMessage(sumoName, endpoint, scanner.Bytes())
	}

	wg.Wait()

	log.Println("processed ", messageCount, " messages.")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func logError(err error, msg []byte) {
	log.Println("Failed to post message", err.Error(), " message: ", string(msg))
}

func statusOk(status int) bool {
	return status >= 200 && status < 300
}

func postMessage(sumoName, endpoint string, msg []byte) {
	defer wg.Done()
	messageCount += 1
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(msg))
	if err != nil {
		logError(err, msg)
		return
	}

	if sumoName != "" {
		req.Header.Add("X-Sumo-Name", sumoName)
	}

	resp, err := client.Do(req)
	if err != nil {
		logError(err, msg)
		return
	}

	defer resp.Body.Close()
	respMsg, err := ioutil.ReadAll(resp.Body)
	if !statusOk(resp.StatusCode) {
		if err == nil {
			log.Println("Message post error: ", respMsg, " status: ", resp.Status)
		} else {
			log.Println("Message post error - status: ", resp.Status)
		}
	}
}
