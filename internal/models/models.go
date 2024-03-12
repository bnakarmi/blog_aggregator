package models

import (
	"encoding/xml"
	"time"

	"github.com/bnakarmi/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CreateUserRequest struct {
	User string `json:"user"`
}

type CreateFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt *time.Time
}

type CreateFeedResponse struct {
	Feed       Feed                `json:"feed"`
	FeedFollow database.FeedFollow `json:"feed_follow"`
}

type CreateFeedFollowsRequest struct {
	FeedID uuid.UUID `json:"feed_id"`
}

type Rss struct {
	XMLName xml.Name   `xml:"rss"`
	Text    string     `xml:",chardata"`
	Atom    string     `xml:"atom,attr"`
	Version string     `xml:"version,attr"`
	Channel RssChannel `xml:"channel"`
}

type RssChannel struct {
	Text          string    `xml:",chardata"`
	Title         string    `xml:"title"`
	Link          RssLink   `xml:"link"`
	Description   string    `xml:"description"`
	Generator     string    `xml:"generator"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Item          []RssItem `xml:"item"`
}

type RssLink struct {
	Text string `xml:",chardata"`
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type RssItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
}

type PostResponse struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt *time.Time
	FeedID      uuid.UUID
}
