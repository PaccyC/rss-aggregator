package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/PaccyC/rss-aggregator/internal/database"
)


func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
){
	log.Printf("Scraping on %v goroutines every %v duration", concurrency, timeBetweenRequest)

	ticker :=time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C{
	 feeds , err :=	db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
	)
	if err != nil {
		log.Println("Could not fetch feed: ",err)
		continue
	}

	wg := &sync.WaitGroup{}
	 
	for _, feed := range feeds {
		wg.Add(1)
    go scrapeFeed(db,wg,feed)
	}
	wg.Wait()
	}
}
func scrapeFeed (db *database.Queries,wg *sync.WaitGroup, feed database.Feed){
 defer	wg.Done()


 
 
 _, err :=db.MarkFeedAsFetched(context.Background(),
feed.ID)
if err != nil {
log.Println("Error marking feed as fetched, ",err)
return
}

rssFeed, err:=urlToFeed(feed.Url)
if err != nil {
	log.Println("Error fetching as fetched, ",err)
	return
	}
	for _ ,item := range rssFeed.Channel.Item{
		log.Println("Post found , ", item.Title)
	}
	log.Printf("Feed %s collected, %v post found", feed.Name, len(rssFeed.Channel.Item))

}