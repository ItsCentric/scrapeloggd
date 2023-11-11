/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"scrapeloggd/webscraper"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var outputDirectory string
var scrapeCmd = &cobra.Command{
	Use:     "scrape URL -o output/directory",
	Example: "scrapeloggd scrape https://www.backloggd.com/u/Centric -o my/cool/directory",
	Short:   "Scrape your Backloggd profile information",
	Long:    `Accepts a Backloggd URL pointing to a user's page`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(outputDirectory)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatal("Directory does not exist")
			} else {
				log.Fatal(err.Error())
			}
		}
		userPageUrl, err := webscraper.ProcessURL(args[0])
		if err != nil {
			log.Fatal(err.Error())
		}
		games, err := webscraper.ScrapeBackloggd(userPageUrl)
		if err != nil {
			log.Fatal(err.Error())
		}
		splitUrl := strings.Split(args[0], "/u/")
		userName := strings.TrimRight(splitUrl[1], "/")
		currentUnixTime := time.Now().Unix()
		file, err := os.Create(fmt.Sprintf("%s/backloggd-%s-%d.csv", outputDirectory, userName, currentUnixTime))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()
		csvWriter := csv.NewWriter(file)
		csvWriter.Write([]string{"ID", "Title", "Cover", "Rating"})
		for i := 0; i < len(games); i++ {
			game := games[i]
			csvWriter.Write([]string{fmt.Sprint(game.Id), game.Title, game.Cover, fmt.Sprint(game.Rating)})
		}
		defer csvWriter.Flush()
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "Directory to output CSV file with data")
	scrapeCmd.MarkFlagRequired("output")
}
