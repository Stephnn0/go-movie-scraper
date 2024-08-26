package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	URL         = "https://www.goojara.to/watch-movies-genre"
	scraperType = "without"
)

type Movie struct {
	Genre         string
	Title         string
	MovieLink     string
	Year          string
	ImageURL      string
	QualityLevels string
}

func main() {

	var elapsedTime time.Duration

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

	genres := []string{"Action", "Adventure", "Comedy", "Drama", "Sci-Fi", "Horror", "Crime", "Thriller", "Romance", "Fantasy"}

	var wg sync.WaitGroup

	ch := make(chan []Movie, len(genres))

	startTime := time.Now()

	switch scraperType {

	case "with":
		for _, genre := range genres {

			wg.Add(1)
			go scrapeMoviesWithGoroutines(&wg, genre, ch)

		}

		go func() {

			wg.Wait()
			close(ch)
		}()

		elapsedTime = time.Since(startTime)

		for movieList := range ch {

			writeMovieToCSV(writer, movieList)
		}
	case "without":
		var allMovieList [][]Movie
		for _, genre := range genres {

			returnMovieList := scrapeMovieWithoutGoroutines(genre)
			allMovieList = append(allMovieList, returnMovieList)
		}

		elapsedTime = time.Since(startTime)

		for _, movieList := range allMovieList {

			writeMovieToCSV(writer, movieList)
		}

	default:
		fmt.Println("invalid input, please enter with or without ")
	}

	fmt.Println("finished ...")

	fmt.Printf("Time taken: %v\n", elapsedTime)

}

func writeMovieToCSV(writer *csv.Writer, movieList []Movie) {

	for _, movie := range movieList {

		if err := writer.Write([]string{

			movie.Genre,
			movie.Title,
			movie.Year,
			movie.QualityLevels,
			movie.MovieLink,
			movie.ImageURL,
		}); err != nil {

			fmt.Println("Error writing data to csv")
		}
	}

}

func scrapeMoviesWithGoroutines(wg *sync.WaitGroup, genre string, ch chan<- []Movie) {

	defer wg.Done()

	c := colly.NewCollector()

	movieList := []Movie{}

	newURL := fmt.Sprintf("%s-%s", URL, genre)

	c.OnHTML("div.dflex", func(e *colly.HTMLElement) {

		e.ForEach("div > a", func(_ int, a *colly.HTMLElement) {

			var movie Movie
			movie.Genre = genre
			movie.Title = a.ChildText("span.mtl")
			movie.MovieLink = a.Attr("href")
			movie.ImageURL = a.ChildAttr("img", "src")
			movie.QualityLevels = a.ChildText("span.hda")
			movie.Year = a.ChildText("span.hdy")
			movieList = append(movieList, movie)

		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting: ", r.URL)
	})

	err := c.Visit(newURL)

	if err != nil {

		log.Fatal(err)

	}

	ch <- movieList

}

func scrapeMovieWithoutGoroutines(genre string) []Movie {

	c := colly.NewCollector()

	movieList := []Movie{}

	newURL := fmt.Sprintf("%s-%s", URL, genre)

	c.OnHTML("div.dflex", func(e *colly.HTMLElement) {
		e.ForEach("div > a", func(_ int, a *colly.HTMLElement) {
			var movie Movie
			movie.Genre = genre
			movie.Title = a.ChildText("span.mtl")
			movie.MovieLink = a.Attr("href")
			movie.ImageURL = a.ChildAttr("img", "src")
			movie.QualityLevels = a.ChildText("span.hda")
			movie.Year = a.ChildText("span.hdy")
			movieList = append(movieList, movie)

		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting: ", r.URL)
	})

	err := c.Visit(newURL)

	if err != nil {

		log.Fatal(err)

	}

	return movieList

}
