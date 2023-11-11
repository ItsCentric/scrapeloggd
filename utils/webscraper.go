package utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Game struct {
	id     int64
	title  string
	rating int64
	cover  string
}

const GAMES_PER_PAGE float64 = 50.0

func ScrapeBackloggd(userPage string) ([]Game, error) {
	var games []Game
	pages, err := getGamePages(userPage + "/games")
	if err != nil {
		return []Game{}, errors.New(err.Error())
	}
	urlMatched, err := regexp.MatchString(`^(https?:\/\/)?(www\.)?backloggd\.com\/u\/[a-zA-Z0-9_-]+\/?$`, userPage)
	if err != nil {
		return []Game{}, errors.New("Could not execute URL regex on URL")
	}
	if urlMatched == false {
		return []Game{}, errors.New("Invalid user URL")
	}
	collector := colly.NewCollector(
		colly.AllowedDomains("backloggd.com", "www.backloggd.com"),
	)

	collector.OnError(handleError)
	var gameError error
	collector.OnHTML("div#game-lists", func(gameListContainer *colly.HTMLElement) {
		gameListContainer.ForEach("div > div[game_id]", func(index int, gameElement *colly.HTMLElement) {
			game, err := parseGame(gameElement)
			if err != nil {
				gameError = err
			}
			fmt.Printf("------------\n%s\n- ID: %d\n- Rating: %d\n- Cover: %s\n", game.title, game.id, game.rating, game.cover)
			games = append(games, game)
		})
	})
	if gameError != nil {
		return []Game{}, errors.New(gameError.Error())
	}

	for i := 1; i <= pages; i++ {
		fmt.Printf("------------\n(%d out of %d) Scraping user's games...\n", i, pages)
		collector.Visit(fmt.Sprintf("%s/games?page=%d", userPage, i))
	}

	return games, nil
}

func getGamePages(userUrl string) (int, error) {
	pages := 0
	var err error
	collector := colly.NewCollector(
		colly.AllowedDomains("backloggd.com", "www.backloggd.com"),
	)

	collector.OnHTML("p.subtitle-text", func(numberOfGamesElement *colly.HTMLElement) {
		elementText := numberOfGamesElement.Text
		separatedElementText := strings.Split(elementText, " ")
		numberOfGames, err := strconv.ParseFloat(separatedElementText[0], 64)
		if err != nil {
			err = errors.New("Couldn't parse number of games: " + err.Error())
			return
		}
		pages = int(math.Ceil(numberOfGames / GAMES_PER_PAGE))
	})

	collector.Visit(userUrl)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	return pages, nil
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

func parseGame(element *colly.HTMLElement) (Game, error) {
	id := element.Attr("game_id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return Game{}, errors.New(fmt.Sprintf("Couldn't parse game ID with ID of %s: %s", id, err.Error()))
	}
	title := element.ChildText("div.game-text-centered")
	rating := element.Attr("data-rating")
	if rating == "" {
		rating = "-1"
	}
	ratingInt, err := strconv.ParseInt(rating, 10, 32)
	if err != nil {
		return Game{}, errors.New(fmt.Sprintf("Couldn't parse game rating with ID of %s: %s", id, err.Error()))
	}
	cover := element.ChildAttr("div > img.card-img", "src")

	return Game{
		id:     idInt,
		title:  title,
		rating: ratingInt,
		cover:  cover,
	}, nil
}

func handleError(response *colly.Response, error error) {
	log.Fatal(fmt.Sprintf("An error occured while scraping:\n%s", error.Error()))
}
