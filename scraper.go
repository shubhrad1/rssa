package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shubhrad1/rssagg/internal/database"
)

func timeParser(timestring string) (time.Time, error) {
	time_format := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC850,
		time.ANSIC,
		time.RFC3339,
	}
	var t time.Time
	var err error
	for _, formats := range time_format {
		t, err = time.Parse(formats, timestring)
		if err == nil {
			return t, err
		}
	}
	return t, err

}

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {

	log.Printf("Scraping on %v go-routines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("[ERROR] Error fetching fieeds:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}

}
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {

	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("[ERROR] Error marking feed as fetched: ", err)
		return
	}
	rssfeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("[ERROR] Error fetching feed: ", err)
		return
	}
	for _, item := range rssfeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		//pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		pubAt, err := timeParser(item.PubDate)
		if err != nil {
			log.Printf("[ERROR] Couldn't parse date %v with error: %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{

			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("[ERROR] Failed to create post: ", err)
		}
	}
	log.Printf("Feeds %s collected, %v posts found", feed.Name, len(rssfeed.Channel.Item))
}
