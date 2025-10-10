package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Scraper interface {
	Scrape() ([]NewsItem, error)
}

type NewsItem struct {
	Title   string
	Author  string
	Summary string
	URL     string
	Date    time.Time
}

type IJFScraper struct{}

func (s IJFScraper) Scrape() ([]NewsItem, error) {
	news := make([]NewsItem, 0)
	maxNews := 15
	var count int

	c := colly.NewCollector(
		colly.AllowedDomains("www.ijf.org", "ijf.org"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*ijf.org*",
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	articleCollector := c.Clone()

	var mu sync.Mutex

	c.OnHTML("a.hero, a.news-item", func(e *colly.HTMLElement) {
		mu.Lock()
		if count >= maxNews {
			mu.Unlock()
			return
		}
		count++
		mu.Unlock()

		link := e.Attr("href")
		fmt.Printf("News found: %q -> %s\n", e.Text, link)
		articleCollector.Visit(e.Request.AbsoluteURL(link))
	})

	articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
		subtitle := e.ChildText("div.subtitle")
		if subtitle == "" {
			log.Printf("No subtitle found for %s", e.Request.URL.String())
		}
		splitSubtitle := strings.SplitN(subtitle, " on ", 2)
		author := strings.TrimPrefix(strings.TrimSpace(splitSubtitle[0]), "By ")
		var parsedDate time.Time
		var err error
		if len(splitSubtitle) == 2 {
			date := strings.TrimSpace(splitSubtitle[1])
			dateToParse := strings.ReplaceAll(date, ".", "")
			parsedDate, err = time.Parse("2 Jan 2006", dateToParse)
			if err != nil {
				log.Printf("Error parsing date: %v", err)
			}
		}

		item := NewsItem{
			Title:   e.ChildText("h1"),
			Author:  author,
			Summary: e.ChildText("div.chunk.chunk--summary"),
			URL:     e.Request.URL.String(),
			Date:    parsedDate,
		}
		mu.Lock()
		news = append(news, item)
		mu.Unlock()
		fmt.Printf("Article scraped: %+v\n", item)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})
	articleCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})

	c.Visit("https://www.ijf.org/news")
	c.Wait()
	articleCollector.Wait()
	return news, nil
}

type WTScraper struct{}

func (s WTScraper) Scrape() ([]NewsItem, error) {
	news := make([]NewsItem, 0)
	maxNews := 15
	var count int

	c := colly.NewCollector(
		colly.AllowedDomains("www.worldtaekwondo.org", "worldtaekwondo.org"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*worldtaekwondo.org*",
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	articleCollector := c.Clone()

	var mu sync.Mutex

	c.OnHTML("ul.news_list li", func(e *colly.HTMLElement) {
		mu.Lock()
		if count >= maxNews {
			mu.Unlock()
			return
		}
		count++
		mu.Unlock()

		subj := e.ChildText("span.subj")
		link := e.ChildAttr("a", "href")
		fmt.Printf("News found: %q -> %s\n", subj, link)
		articleCollector.Visit(e.Request.AbsoluteURL(link))
	})

	articleCollector.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("div.head_title")

		var firstParagraph string
		e.ForEach("div.intend p", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" && firstParagraph == "" {
				firstParagraph = text
			}
		})

		var parsedDate time.Time
		var err error
		re := regexp.MustCompile(`\(([^)]+)\)`)
		match := re.FindStringSubmatch(firstParagraph)
		if len(match) > 1 {
			dateStr := match[1] // e.g. "Aug. 31, 2025"
			// WT uses abbreviated month names with a dot, so remove it
			dateToParse := strings.ReplaceAll(dateStr, ".", "")
			parsedDate, err = time.Parse("Jan 2, 2006", dateToParse)
			if err != nil {
				log.Printf("Error parsing date at %s: %v", e.Request.URL.String(), err)
			}
		}

		item := NewsItem{
			Title:   title,
			Author:  "World Taekwondo",
			Summary: firstParagraph,
			URL:     e.Request.URL.String(),
			Date:    parsedDate,
		}

		mu.Lock()
		news = append(news, item)
		mu.Unlock()
		fmt.Printf("Article scraped: %+v\n", item)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})
	articleCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})

	c.Visit("https://www.worldtaekwondo.org/wtnews/list.html?mcd=C02")
	c.Wait()
	articleCollector.Wait()
	return news, nil
}
