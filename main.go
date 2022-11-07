package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var client *http.Client
var okResponses int = 0

func MakeRequest(url string, ch chan<- string) {
	// postData := bytes.NewBuffer([]byte(`{"post":"Amado Amado Amado Amado Amado"}`))
	// req, _ := http.NewRequest("GET", url, postData)
	start := time.Now()
	resp, err := client.Get(url)
	var status int
	if err == nil {
		status = resp.StatusCode
		if resp.StatusCode == 200 {
			okResponses++
		}
		secs := time.Since(start).Seconds()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			return
			// os.Exit(1)
		}
		ch <- fmt.Sprintf("%.4f seconds elapsed with response length: %d %s %d", secs, len(body), url, status)
		defer resp.Body.Close()
	}
	// resp, _ := client.Do(req)

}

// PRIMEIRO ARGUMENTO SE REFERE AO NÚMERO DE REQUISIÇÕES QUE VAI SER ENVIADO POR SEGUDO POR CADA THREAD
// O SEGUNDO É O TOTAL DE REQUISIÇÕES QUE SERÃO ENVIADAS
func main() {
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 1000, //15,04
		IdleConnTimeout:     60 * time.Second,
	}
	client = &http.Client{Transport: tr}
	// usar função sleep
	start := time.Now() //131
	numReqs, _ := strconv.Atoi(os.Args[1])
	maxReqs, _ := strconv.Atoi(os.Args[2])
	url := "http://php-apache.mathec.com.br/" //"https://dasdga1231289371289kjfhaskljdh.herokuapp.com/" //"https://dasdga1231289371289kjfhaskljdh.herokuapp.com/" //"https://httpbin.org/image/png" // "https://httpbin.org/stream/10"
	ch := make(chan string)
	reqs := 0
	//14.72
	// totalTime := time.Now()
	for {
		for i := 0; i <= runtime.NumCPU(); i++ {
			go func() {
				for i := 0; i < numReqs; i++ {
					if reqs < maxReqs { // note: usar tempo total ao invés de número de iterações no futuro.
						go MakeRequest(url, ch)
						reqs++
					} else {
						break
					}
				}
			}()

			if reqs >= maxReqs { // note: usar tempo total ao invés de número de iterações no futuro.
				break
			}
		}
		// for i := 0; i < numReqs; i++ {
		// 	reqs++
		// 	go MakeRequest(url, ch)
		// }
		time.Sleep(1 * time.Second)
		// if time.Since(totalTime).Minutes() >= 1 { // note: usar tempo total ao invés de número de iterações no futuro.
		// 	break
		// }
		if reqs >= maxReqs { // note: usar tempo total ao invés de número de iterações no futuro.
			break
		}
	}
	for i := 0; i < reqs; i++ {
		fmt.Println(<-ch, " \n ", reqs)
	}

	fmt.Printf("%.4fs elapsed\n", time.Since(start).Seconds())
	println("Responses with 200-ok status:", okResponses, "Total reqs:", reqs)
}
