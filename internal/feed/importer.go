package feed

import (
	"context"
	"fmt"
	"log"

	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/ehumba/kungfu-hub/internal/scraper"
)

type Importer struct {
	db database.Queries
}

func NewImporter(db database.Queries) *Importer {
	return &Importer{db: db}
}

func (i *Importer) ImportIJFNews(ctx context.Context, martialArtID int32) error {
	log.Println("Starting IJF import...")
	newsItems, err := scraper.IJFScraper{}.Scrape()
	if err != nil {
		return fmt.Errorf("failed to scrape IJF news: %w", err)
	}
	eventItems, err := scraper.IJFEventScraper{}.Scrape()
	if err != nil {
		return fmt.Errorf("failed to scrape IJF events: %w", err)
	}

	i.insertNewsItems(ctx, martialArtID, newsItems)
	i.insertEvents(ctx, martialArtID, eventItems)
	log.Println("IJF import completed.")
	return nil
}

func (i *Importer) ImportWTNews(ctx context.Context, martialArtID int32) error {
	log.Println("Starting WT import...")
	newsItems, err := scraper.WTScraper{}.Scrape()
	if err != nil {
		return fmt.Errorf("failed to scrape WT news: %w", err)
	}
	eventItems, err := scraper.WTEventScraper{}.Scrape()
	if err != nil {
		return fmt.Errorf("failed to scrape WT events: %w", err)
	}

	i.insertNewsItems(ctx, martialArtID, newsItems)
	i.insertEvents(ctx, martialArtID, eventItems)
	log.Println("WT import completed.")
	return nil
}

func (i *Importer) ImportIBJJFNews(ctx context.Context, martialArtID int32) error {
	log.Println("Starting IBJJF import...")
	eventItems, err := scraper.IBJJFEventScraper{}.Scrape()
	if err != nil {
		return fmt.Errorf("failed to scrape IBJJF events: %w", err)
	}

	i.insertEvents(ctx, martialArtID, eventItems)
	log.Println("IBJJF import completed.")
	return nil
}
