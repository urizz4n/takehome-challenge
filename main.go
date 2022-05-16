package main

import (
	"fmt"
	"sync"
	"takehome-challenge/downloader"
	"takehome-challenge/fetcher"
	"takehome-challenge/houseresponse"
	"time"
)

const (
	url        = "https://app-homevision-staging.herokuapp.com/api_project/houses"
	totalpages = 10
	maxRetries = 5 // Max retries for backoff
)

func execute() {
	var wg sync.WaitGroup

	for i := 1; i <= totalpages; i++ {
		wg.Add(1)

		// We pass the "i" value so it knows which page it is
		go func(i int) {

			var result houseresponse.HouseResponse

			// Checks the page we attempt to fetch is reachable
			fetcher.Fetch(url, maxRetries, i, &result)

			// Downloads all the houses from the page
			downloader.PrepareDownload(maxRetries, &result, &wg)

			fmt.Println("Finished downloading page:", i)
			wg.Done()

		}(i)
	}
	wg.Wait()

	fmt.Printf("Finished downloading all photos for the first %d pages\n", totalpages)
}

func main() {

	start := time.Now()
	execute()
	elapsed := time.Since(start)
	fmt.Printf("Total execution time: %v", elapsed)

}
