# Kungfu Hub
A martial arts news and events aggregator that helps practitioners stay updated with their favorite martial arts disciplines.

## Features
- **News and Event Aggregation:** Automatically scrapes and aggregates news and events from major martial arts organizations:
  - [International Judo Federation (IJF)](https://www.ijf.org)
  - [World Taekwondo (WT)](https://www.worldtaekwondo.org)
  - [International Brazilian Jiu-Jitsu Federation (IBJJF)](https://ibjjf.com)
- **Personalized Feed:** Subscribe to specific martial arts and view a custom-tailored feed.
- **Social Features:** Follow other practitioners and see their comments.
- **JWT Authentication:** Secure login and token-based authentication.

## Built With
- [Go](https://go.dev) — Backend API and scrapers
- [PostgreSQL](https://www.postgresql.org) — Database
- [Docker Compose](https://docs.docker.com/compose/) — Container orchestration
- [Goose](https://github.com/pressly/goose) — Database migrations
- [SQLC](https://sqlc.dev) — Type-safe database access
- [Colly](https://github.com/gocolly/colly) — Web scraping


## Quick Start
### Prerequisites
- Docker and Docker Compose
- Go 1.21 or later
- PostgreSQL 15
- Make (optional)
### Installation
1. Clone the repository
```
git clone https://github.com/ehumba/kungfu-hub.git
cd kungfu-hub
```
2. Set up environment variables
```
cp .env.example .env
# Edit .env with your preferred settings
```
3. Start the database
`docker-compose up -d db`
4. Run database migrations
```
make migrate-up
# Or without make:
goose -dir ./sql/schema postgres "<db_connection_string>" up
```
5. Start the server
`go run .`

The API will be available at `http://localhost:8080`

## API Endpoints
### Authentication
- `POST /api/users` - Create new user
- `POST /api/login` - Login and receive access token
- `POST /api/refresh` - Refresh JWT token
- `DELETE /api/revoke` - Revoke token
### Martial Arts & Content
- `GET /api/martial_arts` - List available martial arts
- `GET /api/events` - List all events
- `GET /api/events_by_art` - List events for specific martial art
- `GET /api/news_by_art` - List all news
- `GET /api/news_by_art` - List news for specific martial art
- `GET /api/feed` - Get personalized feed
### User Actions
- `POST /api/subscribe` - Subscribe to martial art
- `POST /api/unsubscribe` - Unsubscribe from martial art
- `GET /api/subscriptions` - List user's subscriptions
- `POST /api/comments` - Post a comment 
- `POST /api/follow` - Follow another user
- `POST /api/unfollow` - Unfollow another user


## Development
### Database Migrations
```
# Create new migration
goose -dir ./sql/schema create name_of_migration sql

# Run migrations
make migrate-up

# Rollback
make migrate-down
```
### Generate Database Code
```
# After modifying SQL queries
sqlc generate
```
### Data Import
To manually trigger data imports from all martial arts organizations:
`POST /api/admin/import`

## Possible Future Improvements
- Add more supported martial arts
- Add IBJJF news scraping support
- Add UI front-end for browsing feeds
- Integrate ranking or event results data

## License
MIT

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
