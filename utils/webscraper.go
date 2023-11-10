package utils

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func ScrapeBackloggd(userPage string) {
	urlMatched, err := regexp.MatchString(`^(https?:\/\/)?(www\.)?backloggd\.com\/u\/[a-zA-Z0-9_-]+\/$`, userPage)
	if err != nil {
		log.Fatal("Could not execute regex")
		return
	}
	if urlMatched == false {
		log.Fatal("Invalid URL")
		return
	}
	collector := colly.NewCollector(
		colly.AllowedDomains("backloggd.com", "www.backloggd.com"),
	)

	collector.OnResponse(handleResponse)
	collector.OnError(handleError)
	collector.OnHTML("h3.main-header", handleHTML)

	collector.Visit(userPage)
}

func ProcessURL(url string) (string, error) {
	protocolMatched, err := regexp.MatchString(`^(https?://)`, url)
	if err != nil {
		return "", errors.New("Failed to execute protocol regex on URL string")
	}
	if protocolMatched == false {
		url = "https://" + url
	}

	return url, nil
}

func handleHTML(element *colly.HTMLElement) {
	fmt.Println("Your username:", strings.TrimSpace(element.Text))
}

func handleResponse(response *colly.Response) {
	fmt.Println("Scraping your profile...")
}

func handleError(response *colly.Response, error error) {
	log.Fatal(fmt.Sprintf("An error occured while scraping:\n%s", error.Error()))
}
