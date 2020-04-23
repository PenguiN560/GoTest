package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func bodyCounter(url string) int {
	var (
		client http.Client
		count  int = 0
	)
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		count = strings.Count(bodyString, "Go")

		fmt.Printf("Count for %s: %d\n", url, count)
	}

	return count
}

func counter(urlsChannel <-chan string, countersChannel chan<- int) {
	var count int = 0
	for {
		url, more := <-urlsChannel
		if more {
			count += bodyCounter(url)
		} else {
			countersChannel <- count
			return
		}
	}
}

func main() {
	const k = 5
	var (
		i     int = 0
		count int = 0
	)

	urlsChannel := make(chan string, k)
	countersChannel := make(chan int, k)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urlsChannel <- scanner.Text()
		if i < k {
			go counter(urlsChannel, countersChannel)
			i++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(urlsChannel)
	for ; i > 0; i-- {
		count += <-countersChannel
	}

	fmt.Println("Total: ", count)
}
