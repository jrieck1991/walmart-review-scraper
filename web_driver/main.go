package main

import (
	"sync"
	"encoding/csv"
	"strings"
	"fmt"
	"os"
	"github.com/tebeka/selenium"
)

const (
	seleniumPath = "./drivers/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath = "./drivers/geckodriver"
	port = 8080

	walmartURL = "https://www.walmart.com/reviews/product/745013883"
	reviews = "div.Grid.ReviewList-content"
	nextBtn = "button.paginator-btn.paginator-btn-next"
	csvFilePath = "/go/src/web_driver/walmart_reviews.csv"
)

func main() {
	var err error

	// selenium options
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // specify the path to GeckoDriver in order to use firefox.
		selenium.Output(os.Stderr), // output debug information to stderr
	}
	
	// start selenium
	svc, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer svc.Stop()

	// Connect to webdriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// navigate to target url
	if err = wd.Get(walmartURL); err != nil {
		panic(err)
	}

	// create csv file
	f, err := os.Create(csvFilePath)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	// init csv writer & write headers to csv file
	w := csv.NewWriter(file)
	if err := w.Write([]string{"title", "rating", "comment", "userDate", "domain", "overflow"}); err != nil {
		panic(err)
	}
	

	// first page elements
	var e []selenium.WebElement
	e, err = wd.FindElements(selenium.ByCSSSelector, reviews)
	if err != nil {
		panic(err)
	}

	// write elements to csv file
	for _, s := range e {
		// get text
		t, err := s.Text()
		if err != nil {
			panic(err)
		}

		// write to csv
		r := strings.Split(t, "\n")
		if err := w.Write(r); err != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup
	for {
		// Get reviews on a page
		n, err := wd.FindElements(selenium.ByCSSSelector, reviews)
		if err != nil {
			fmt.Println(err)
			break
		}

		// write reviews to csv file
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for _, s := range n {
				// get text
				t, err := s.Text()
				if err != nil {
					fmt.Println(err)
					break
				}
	
				// write to csv
				r := strings.Split(t, "\n")
				if err := w.Write(r); err != nil {
					fmt.Println(err)
					break
				}
			}
		}(&wg)

		// go to next page
		btn, err := wd.FindElement(selenium.ByCSSSelector, nextBtn)
		if err != nil {
			fmt.Println(err)
			break
		}
		if err := btn.Click(); err != nil {
			fmt.Println(err)
			break
		}

	}
	wg.Wait()
}