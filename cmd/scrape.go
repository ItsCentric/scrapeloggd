/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"scrapeloggd/webscraper"

	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:     "scrape URL",
	Example: "scrapeloggd scrape https://www.backloggd.com/u/Centric",
	Short:   "Scrape your Backloggd profile information",
	Long:    `Accepts a Backloggd URL pointing to a user's page`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userPageUrl, err := webscraper.ProcessURL(args[0])
		if err != nil {
			log.Fatal(err.Error())
		}
		webscraper.ScrapeBackloggd(userPageUrl)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
