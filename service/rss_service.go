package service

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bnakarmi/blog_aggregator/internal/models"
	"github.com/bnakarmi/blog_aggregator/repository"
)

type FeedWorker struct {
	Repository *repository.QueryRepository
}

func NewFeedWorkerService(repository *repository.QueryRepository) *FeedWorker {
	return &FeedWorker{
		Repository: repository,
	}
}

func (worker *FeedWorker) Initialize() {
	repository := worker.Repository
	context := context.Background()
	ticker := time.NewTicker(60 * time.Second)
	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	// Start the worker in the background
	go func() {
		for t := range ticker.C {
			log.Println("Feed fetching started at", t)

			nextFeeds, err := repository.GetNextFeedsToFetch(context)
			if err != nil {
				log.Printf("Unable to get feeds %v", err)
				continue
			}

			if len(nextFeeds) == 0 {
				log.Println("No feeds found")
				continue
			}

			for _, feed := range nextFeeds {
				wg.Add(1)
				semaphore <- struct{}{}

				go func() {
					defer wg.Done()
					defer func() { <-semaphore }()

					rssFeed, err := fetchFeeds(feed.Url)
					if err != nil {
						log.Printf("Unable to get feeds %v", err)
						return
					}

					for _, rssFeedItem := range rssFeed.Channel.Item {
                        _, err := repository.CreatePost(context, &rssFeedItem, feed.ID)
                        if err != nil {
                            log.Printf("Error creating post %v", err)
                        }
					}
				}()
			}
		}
	}()
}

func fetchFeeds(url string) (models.Rss, error) {
	log.Println("Feed url:", url)
	var rssFeeds models.Rss

	resp, err := http.Get(url)
	if err != nil {
		return rssFeeds, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rssFeeds, err
	}

	err = xml.Unmarshal(body, &rssFeeds)
	if err != nil {
		return rssFeeds, nil
	}

	return rssFeeds, nil
}
