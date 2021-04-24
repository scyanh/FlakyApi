package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type HouseResponse struct {
	Houses []House `json:"houses"`
}
type House struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

// create the filename with the format: id-[NNN]-[address].[ext]
func (h *House) GetFilename() string {
	// replace spaces and punctuations
	address := strings.ReplaceAll(h.Address, " ", "_")
	address = strings.ReplaceAll(address, ".", "")
	address = strings.ReplaceAll(address, ",", "")

	// split the PhotoURL to get the ext
	photoSplit := strings.Split(h.PhotoURL, ".")

	// build the filename
	filename := fmt.Sprintf("id-%d-%s.%s", h.Id, address, photoSplit[len(photoSplit)-1])
	return filename
}

// download file
func (h *House) DownloadFile() error {
	//Get the response bytes from the url
	response, err := http.Get(h.PhotoURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	//Create a empty file
	file, err := os.Create(h.GetFilename())
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
