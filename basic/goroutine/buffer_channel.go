package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	work := make(chan string, 1024)
	numWorker := 10

	var wg sync.WaitGroup

	for i := 0; i < numWorker; i++ {
		go webgetworker(work, &wg)
	}

	urls := []string{
		"http://example.com",
		"http://www.baidu.com",
		"http://www.163.com",
		"http://www.qq.com",
		"http://www.360.cn",
	}

	for i := 0; i < 100; i++ {
		for _, url := range urls {
			wg.Add(1)
			work <- url
		}
	}

	wg.Wait()
}

func webgetworker(in <-chan string, wg *sync.WaitGroup) {
	for {
		url := <-in

		res, err := http.Get(url)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Get %s : %d\n", url, res.StatusCode)
		}

		wg.Done()

	}
}
