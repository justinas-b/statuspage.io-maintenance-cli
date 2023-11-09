package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/justinas-b/statuspage.io-maintenance-cli/client"
	tuiConfirmation "github.com/justinas-b/statuspage.io-maintenance-cli/tui/confirmation"
	tuiList "github.com/justinas-b/statuspage.io-maintenance-cli/tui/list"
	tuiProgress "github.com/justinas-b/statuspage.io-maintenance-cli/tui/progress"
	tuiTextarea "github.com/justinas-b/statuspage.io-maintenance-cli/tui/textarea"
	tuiTextinput "github.com/justinas-b/statuspage.io-maintenance-cli/tui/textinput"
	"os"
	"strings"
)

func init() {
}

var (
	apiKeys string
	pages   []*client.Page
)

func getAPIKeysFromArgs() error {
	flag.StringVar(&apiKeys, "apiKeys", "", "Comma-separated list of API keys")
	flag.Parse()
	if f := flag.CommandLine.Lookup("apiKeys"); f.Value.String() == "" {
		return errors.New("apiKeys argument not set")
	}
	return nil
}

func getAPIKeysFromEnv() error {
	var found bool
	apiKeys, found = os.LookupEnv("STATUSPAGE_API_KEYS")
	if !found {
		return errors.New("STATUSPAGE_API_KEYSenvironment variable not set")
	}
	return nil
}

func main() {
	err := getAPIKeysFromArgs()
	if err != nil {
		err := getAPIKeysFromEnv()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	for _, apiKey := range strings.Split(apiKeys, ",") {
		c := client.NewClient(apiKey)
		p, _ := c.GetPages()
		pages = append(pages, p...)
	}

	filteredPages := tuiList.Run(pages)

	maintenanceTitle := tuiTextinput.Run(
		"Please enter title for new maintenance",
		"Maintenance title",
		"MAINTENANCE TITLE CANNOT BE EMPTY",
		validateEmptyString,
	)

	maintenanceDescription := tuiTextarea.Run()

	maintenanceStartDate := tuiTextinput.Run(
		"Please enter start date in UTC time zone",
		"YYYY-MM-DD HH:MM:SS, for example 2023-01-31 20:30:00",
		"Provided format does not match YYYY-MM-DD HH:MM:SS or is not in the future",
		validateDateTimeFormat,
	)

	maintenanceDuration := tuiTextinput.Run(
		"Please enter maintenance duration in hours",
		"Integer, for example 8",
		"Provided value is not a valid integer",
		validateInteger,
	)

	confirmation := tuiConfirmation.Run(
		fmt.Sprintf(
			"Please confirm if all details below are correct?\n\n"+
				"Maintenance title:\t\t%s\n"+
				"Maintenance description:\t%s\n"+
				"Maintenance start date:\t\t%s\n"+
				"Maintenance duration:\t\t%s\n"+
				"Status pages:\t\t\t%+q\n"+
				"\n",
			maintenanceTitle, maintenanceDescription, maintenanceStartDate, maintenanceDuration, filteredPages),
		[]string{"Yes", "No"},
	)

	if confirmation != "Yes" {
		fmt.Printf("Aborted\n")
		os.Exit(1)
	}

	tuiProgress.Run(filteredPages, maintenanceTitle, maintenanceDescription, maintenanceStartDate, maintenanceDuration)

}
