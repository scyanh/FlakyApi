package main

import (
	"fmt"
	"github.com/scyanh/FlakyApi/services"
)

func main() {
	urlPath := "http://app-homevision-staging.herokuapp.com/api_project/houses?page=1"
	apiService := services.NewApiService()
	houses, err := apiService.RequestAPI(urlPath)
	if err != nil {
		fmt.Println("err=", err)
		return
	}

	apiService.DownloadFiles(houses)
}
