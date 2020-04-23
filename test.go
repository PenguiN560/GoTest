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

func counter(urlsChannel <-chan string, countersCannel chan<- int) {
	var counter int = 0
	for {
		url, more := <-urlsChannel
		if more {
			var client http.Client
			resp, err := client.Get(url)
			if err != nil {
				log.Fatal(err)
			}

			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				bodyString := string(bodyBytes)
				count := strings.Count(bodyString, "Go")

				fmt.Printf("Count for %s: %d\n", url, count)
				counter += count
			}

			err = resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			countersCannel <- counter
			return
		}
	}
}

func main() {
	const k = 5
	var (
		i int = 0
		count int = 0
	)

	urlsChannel := make(chan string, k)
	countersCannel := make(chan int, k)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		urlsChannel <- scanner.Text()
		if (i < k) {
			go counter(urlsChannel, countersCannel)
			i++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(urlsChannel)
	for ; i > 0; i-- {
		count += <-countersCannel
	}

	fmt.Println("Total: ", count)
}
