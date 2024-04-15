package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	q := "Nigeria"
	if len(os.Args) > 1 {
		q = os.Args[1]
	}

	// get environment variables
	url := fmt.Sprintf("%v=%v", os.Getenv("URL"), q)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatal("api not reachable!")
	}

	resp, _ := io.ReadAll(res.Body)
	var weather Weather
	if err = json.Unmarshal(resp, &weather); err != nil {
		log.Fatal("unable to unmarshal json")
	}
	location, hours, current := weather.Location, weather.Forecast.Forecastday[0].Hour, weather.Current

	fmt.Printf(
		"%s, %s: %0.fC - %s \n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)
	for _, h := range hours {
		date := time.Unix(h.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		message := fmt.Sprintf(
			"%s - %0.fC, %0.f%%, %s\n",
			date.Format("15:04"),
			h.TempC,
			h.ChanceOfRain,
			h.Condition.Text,
		)

		if h.ChanceOfRain < 40 {
			color.Green(message)
		} else {
			color.Red(message)
		}
	}
}
