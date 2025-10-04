package scraper

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type EventScraper interface {
	Scrape() ([]EventItem, error)
}

type EventItem struct {
	Title       string
	Location    string
	Date        time.Time
	URL         string
	Description string
}

type IJFEventScraper struct{}

func (s IJFEventScraper) Scrape() ([]EventItem, error) {
	events := make([]EventItem, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.ijf.org", "ijf.org"),
		colly.Async(true),
	)

	eventCollector := c.Clone()

	var mu sync.Mutex

	c.OnHTML("a.event-link-title", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Event found: %q -> %s\n", e.Text, link)
		eventCollector.Visit(e.Request.AbsoluteURL(link))
	})

	eventCollector.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("div.title")
		location := e.ChildText("div.location")
		date := e.ChildText("table.table--2017.table--no-header tr:nth-of-type(1) td:nth-of-type(2)")
		dateToParse := strings.ReplaceAll(date, ".", "")
		parsedDate, err := time.Parse("2 Jan 2006", dateToParse)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
		}

		event := EventItem{
			Title:       title,
			Location:    location,
			Date:        parsedDate,
			Description: "",
			URL:         e.Request.URL.String(),
		}
		mu.Lock()
		events = append(events, event)
		mu.Unlock()
		fmt.Println("Event scraped:", event)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})
	eventCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})

	c.Visit("https://www.ijf.org/calendar?age=all")
	c.Wait()
	eventCollector.Wait()
	return events, nil
}

type WTEventScraper struct{}

func (s WTEventScraper) Scrape() ([]EventItem, error) {
	events := make([]EventItem, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.worldtaekwondo.org", "worldtaekwondo.org"),
		colly.Async(true),
	)

	eventCollector := c.Clone()

	var mu sync.Mutex

	c.OnHTML("span.inform_lst a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Event found: %q -> %s\n", e.Text, link)
		eventCollector.Visit(e.Request.AbsoluteURL(link))
	})

	eventCollector.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("div.head_title")

		event := EventItem{
			Title: title,
		}

		mu.Lock()
		events = append(events, event)
		mu.Unlock()
		fmt.Println("Event scraped:", event)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})
	eventCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})

	c.Visit("https://www.worldtaekwondo.org/competition/competition_index.html")
	c.Wait()
	eventCollector.Wait()
	return events, nil
}

type IWUFEventScraper struct{}

func (s IWUFEventScraper) Scrape() ([]EventItem, error) {
	events := make([]EventItem, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.iwuf.org", "iwuf.org"),
		colly.Async(true),
	)

	var mu sync.Mutex

	c.OnHTML("tr.table__event", func(e *colly.HTMLElement) {
		dateStr := strings.TrimSpace(e.ChildText("td p"))
		title := strings.TrimSpace(e.ChildText("td h4"))
		location := strings.TrimSpace(e.ChildText("td.table__location p"))

		// IWUF date format example: "4.5- 4.10"
		// → we’ll take the first part as the start date, assuming current year
		now := time.Now()
		dateParts := strings.Split(dateStr, "-")
		startDateRaw := strings.TrimSpace(strings.ReplaceAll(dateParts[0], ".", "-"))
		startDate := fmt.Sprintf("%d-%s", now.Year(), startDateRaw)
		parsedDate, err := time.Parse("2006-1-2", startDate)
		if err != nil {
			log.Printf("Error parsing date %q: %v", dateStr, err)
			parsedDate = time.Time{}
		}

		item := EventItem{
			Title:       title,
			Location:    location,
			Date:        parsedDate,
			Description: "",
		}

		mu.Lock()
		events = append(events, item)
		mu.Unlock()

		fmt.Printf("Event scraped: %+v\n", item)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://www.iwuf.org/en/calendar/index.html")
	c.Wait()

	return events, nil
}
