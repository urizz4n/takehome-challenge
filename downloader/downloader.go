package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"takehome-challenge/fetcher"
	"takehome-challenge/houseresponse"
)

func downloadFile(maxRetries int, URL, fileName string) {

	var resp *http.Response
	connected := false

	resp, connected = fetcher.BackOff(maxRetries, URL, resp)

	if connected {
		//Create a empty file
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal("Error when creating the file")
		}
		defer file.Close()

		//Write the bytes to the file
		defer resp.Body.Close()
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Fatal("Error when copying the file")
		}

	}
}

func PrepareDownload(maxRetries int, result *houseresponse.HouseResponse, wg *sync.WaitGroup) {

	for _, house := range result.Houses {

		wg.Add(1)

		go func(house houseresponse.Houses) {

			fileExtension := "." + strings.Split(house.PhotoURL, ".")[len(strings.Split(house.PhotoURL, "."))-1]
			fileName := "images/" + strconv.Itoa(house.Id) + "-" + house.Address + fileExtension
			imageUrl := house.PhotoURL

			downloadFile(maxRetries, imageUrl, fileName)

			wg.Done()

		}(house)
	}
}
