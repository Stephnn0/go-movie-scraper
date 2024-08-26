package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type Movie struct {
	Genre         string
	Title         string
	MovieLink     string
	Year          string
	ImageURL      string
	QualityLevels string
}

const (
	URL         = "https://www.goojara.to/watch-movies-genre"
	scraperType = "without"
)

var elapsedTime time.Duration

func main() {

	file, err := os.Create("new-movies.csv")

	if err != nil {
		fmt.Println("error while creating file", err)
		return

	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"Genre", "Title", "Year", "QualityLevels", "MovieLink", "ImageURL"})

	fmt.Println("finished ...")

}
