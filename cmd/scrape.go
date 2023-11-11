/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	webscrapper "scrapeloggd/utils"

	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:     "scrape URL",
	Example: "scrapeloggd scrape https://www.backloggd.com/u/Centric",
	Short:   "Scrape your Backloggd profile information",
	Long:    `Accepts a Backloggd URL pointing to a user's page`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userPageUrl, err := utils.ProcessURL(args[0])
		if err != nil {
			log.Fatal(err.Error())
		}
		utils.ScrapeBackloggd(userPageUrl)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
