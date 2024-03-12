package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bnakarmi/blog_aggregator/internal/database"
	"github.com/bnakarmi/blog_aggregator/internal/models"
)

const InvalidRequestMsg = "Invalid request"

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error marshalling JSON:", err)
		return
	}

	w.Header().Set(ContentTypeKey, ContentTypeJSON)
	w.WriteHeader(status)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, status int, msg string) {
	RespondWithJSON(w, status, models.Response{Status: status, Message: msg})
}

func DecodeJSONFromBody(w http.ResponseWriter, r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(target)
	if err != nil {
		log.Printf("%s %v", InvalidRequestMsg, err)
		http.Error(w, InvalidRequestMsg, http.StatusBadRequest)
		return err
	}

	return nil
}

func GenerateApiKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	return hex.EncodeToString(bytes)
}

func GetApiKey(r *http.Request) (string, error) {
	apiKeyHeader := r.Header.Get("Authorization")
	if apiKeyHeader == "" {
		return "", errors.New("Missing API key")
	}

	apiKeyTokens := strings.Split(apiKeyHeader, " ")
	if len(apiKeyTokens) != 2 || apiKeyTokens[0] != "ApiKey" {
		return "", errors.New("Invalid authorization")
	}

	return apiKeyTokens[1], nil
}

func GetFeedResponse(feed *database.Feed) models.Feed {
	var lastFetchedAt *time.Time
	if feed.LastFetchedAt.Valid {
		lastFetchedAt = &feed.LastFetchedAt.Time
	}

	return models.Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: lastFetchedAt,
	}
}

func GetPostResponse(post *database.Post) models.PostResponse {
	description := ""
	if post.Description.Valid {
		description = post.Description.String
	}

	var publishedAt *time.Time
	if post.PublishedAt.Valid {
		publishedAt = &post.PublishedAt.Time
	}

	return models.PostResponse{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: description,
		PublishedAt: publishedAt,
		FeedID:      post.FeedID,
	}
}
