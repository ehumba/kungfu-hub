package feed

import (
	"context"
	"database/sql"
	"log"

	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/ehumba/kungfu-hub/internal/scraper"
)

func (i *Importer) insertNewsItems(ctx context.Context, martialArtID int32, items []scraper.NewsItem) {
	for _, item := range items {
		_, err := i.db.InsertNewsItem(ctx, database.InsertNewsItemParams{
			Column1: sql.NullInt32{Valid: false},                     // unknown event_id
			Column2: sql.NullInt32{Int32: martialArtID, Valid: true}, // martial_art_id
			Column3: sql.NullString{String: item.Title, Valid: true}, // title
			Column4: sql.NullString{String: item.Summary, Valid: true},
			Column5: sql.NullTime{Time: item.Date, Valid: true},    // date
			Column6: sql.NullString{String: item.URL, Valid: true}, // url
		})
		if err != nil {
			log.Printf("Failed to insert %q: %v", item.Title, err)
			continue
		}
	}

}

func (i *Importer) insertEvents(ctx context.Context, martialArtID int32, events []scraper.EventItem) {
	for _, event := range events {
		_, err := i.db.InsertEvent(ctx, database.InsertEventParams{
			Column1: sql.NullInt32{Int32: martialArtID, Valid: true},  // martial_art_id
			Column2: sql.NullString{String: event.Title, Valid: true}, // title
			Column3: sql.NullTime{Time: event.Date, Valid: true},      // date
			Column4: sql.NullString{String: event.Location, Valid: true},
			Column5: sql.NullString{String: event.Description, Valid: true},
			Column6: sql.NullString{String: event.URL, Valid: true}, // url
		})
		if err != nil {
			log.Printf("Failed to insert event %q: %v", event.Title, err)
			continue
		}
	}
}
