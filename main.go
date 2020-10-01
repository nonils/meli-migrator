package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const FileUrl = ""
const NickName = ""

func main() {
	getMeliItemsByName(NickName)
}

func getMeliItemsByName(sellerName string) {
	record := make(map[string]interface{})
	client := &http.Client{}
	req, err := createRequestForGetMeliItems(sellerName)
	if err != nil {
		fmt.Errorf("puta madre")
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Errorf("puta madre 2")
	}
	if err := json.NewDecoder(res.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	results, ok := record["results"].([]interface{})
	if !ok {
		errors.New("La cagamos nivel 3")
	}

	for _, result := range results {
		url := result.(map[string]interface{})["permalink"]
		imageUrls := getImages(url.(string))
		for _, imageUrl := range imageUrls {
			ss := strings.Split(imageUrl, ".")
			s := ss[len(ss)-1]
			downloadFile(imageUrl, FileUrl+uuid.Must(uuid.NewRandom()).String()+"."+s)
		}

	}
}

func getImages(url string) []string {
	imageUrl := make([]string, 0)
	c := colly.NewCollector()
	c.OnHTML("[class=ui-pdp-gallery__figure] img", func(e *colly.HTMLElement) {
		imageUrl = append(imageUrl, e.Attr("data-zoom"))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})
	c.Visit(url)
	return imageUrl
}

func createRequestForGetMeliItems(sellerName string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://api.mercadolibre.com/sites/MLA/search?nickname="+sellerName, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, errors.New("la cagamos feo")
	}
	return req, nil
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
