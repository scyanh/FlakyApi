package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/scyanh/FlakyApi/models"
	"io/ioutil"
	"net/http"
	"sync"
)

type ApiService struct{}

func NewApiService() ApiInterface {
	return &ApiService{}
}

type ApiInterface interface {
	RequestAPI(urlPath string) ([]models.House, error)
	DownloadFiles(houses []models.House)
}

// Request to Homevision API
func (a ApiService) RequestAPI(urlPath string) ([]models.House, error) {
	res, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("received non 200 response code")
	}

	body, _ := ioutil.ReadAll(res.Body)
	var apiResponse models.HouseResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.Houses, nil
}

// Download photos with workers concurrently
func (a ApiService) DownloadFiles(houses []models.House) {
	// channel for download images
	jobs := make(chan models.House)

	// start workers
	wg := &sync.WaitGroup{}
	wg.Add(maxWorkers)
	for i := 1; i <= maxWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			for j := range jobs {
				doWork(i, j)
			}
		}(i)
	}

	// add jobs
	for _, house := range houses {
		jobs <- house
	}
	close(jobs)

	// wait for workers to complete
	wg.Wait()

	fmt.Printf("\nAll images downloaded:\n")
}

// max number of workers
const maxWorkers = 4

func doWork(id int, house models.House) {
	fmt.Printf("worker%d: start working for id:%d\n", id, house.Id)
	house.GetFilename()
	err := house.DownloadFile()
	if err != nil {
		fmt.Printf("worker%d: err:%s!\n", id, err)
	} else {
		fmt.Printf("worker%d: completed id:%d!\n", id, house.Id)
	}
}
