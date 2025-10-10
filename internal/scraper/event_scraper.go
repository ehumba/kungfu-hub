package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	maxEvents := 15
	var count int

	c := colly.NewCollector(
		colly.AllowedDomains("www.ijf.org", "ijf.org"),
		colly.Async(true),
	)

	eventCollector := c.Clone()

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*ijf.org*",
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	eventCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*ijf.org*",
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	var mu sync.Mutex

	c.OnHTML("a.event-link-title", func(e *colly.HTMLElement) {
		mu.Lock()
		if count >= maxEvents {
			mu.Unlock()
			return
		}
		count++
		mu.Unlock()

		link := e.Attr("href")
		fmt.Printf("Event found: %q -> %s\n", e.Text, link)
		eventCollector.Visit(e.Request.AbsoluteURL(link))
	})

	eventCollector.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("div.title")
		location := e.ChildText("div.location")
		rawDate := strings.TrimSpace(e.ChildText("table.table--2017.table--no-header tr:nth-of-type(1) td:nth-of-type(2)"))

		// Clean and parse date
		parts := strings.Fields(rawDate)
		var parsedDate time.Time
		if len(parts) >= 3 {
			dateStr := fmt.Sprintf("%s %s %s", parts[0], parts[1], parts[2])
			dateStr = strings.ReplaceAll(dateStr, ".", "") // remove dots in month abbreviation
			var err error
			parsedDate, err = time.Parse("2 Jan 2006", dateStr)
			if err != nil {
				log.Printf("Error parsing date at %s: %v", e.Request.URL.String(), err)
			}
		} else {
			log.Printf("Could not parse date at %s: %q", e.Request.URL.String(), rawDate)
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
	maxEvents := 15
	var count int

	c := colly.NewCollector(
		colly.AllowedDomains("www.worldtaekwondo.org", "worldtaekwondo.org"),
		colly.Async(true),
	)

	eventCollector := c.Clone()

	var mu sync.Mutex

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*worldtaekwondo.org*",
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	c.OnHTML("span.inform_lst a", func(e *colly.HTMLElement) {
		mu.Lock()
		if count >= maxEvents {
			mu.Unlock()
			return
		}
		count++
		mu.Unlock()

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

// IBJJFEventScraper fetches upcoming events from the IBJJF API

type IBJJFEventScraper struct{}

type ibjjfEventItem struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	City            string `json:"city"`
	State           string `json:"state"`
	Country         string `json:"country"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	RegistrationURL string `json:"registration_url"`
}

func (s IBJJFEventScraper) Scrape() ([]EventItem, error) {
	url := "https://ibjjf.com/api/v1/events/upcomings.json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Referer", "https://ibjjf.com/events/championships")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://ibjjf.com")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IBJJF events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Try to decode both array and object formats
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var arr []ibjjfEventItem
	if err := json.Unmarshal(body, &arr); err != nil {
		// fallback: maybe the JSON is an object with a nested array
		var obj map[string]json.RawMessage
		if err2 := json.Unmarshal(body, &obj); err2 == nil {
			for _, v := range obj {
				if err3 := json.Unmarshal(v, &arr); err3 == nil && len(arr) > 0 {
					break
				}
			}
		}
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("no events found in API response")
	}

	maxEvents := 15
	if len(arr) > maxEvents {
		arr = arr[:maxEvents]
	}

	events := make([]EventItem, 0, len(arr))
	for _, ev := range arr {
		start, _ := time.Parse("2006-01-02", ev.StartDate)

		location := ev.City
		if ev.State != "" {
			location += ", " + ev.State
		}
		if ev.Country != "" {
			location += ", " + ev.Country
		}

		events = append(events, EventItem{
			Title:       ev.Name,
			Location:    location,
			Date:        start,
			Description: "",
			URL:         ev.RegistrationURL,
		})
	}

	fmt.Printf("Scraped %d IBJJF events\n", len(events))
	return events, nil
}
