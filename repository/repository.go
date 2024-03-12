package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bnakarmi/blog_aggregator/internal/database"
	"github.com/bnakarmi/blog_aggregator/internal/models"
	"github.com/bnakarmi/blog_aggregator/utils"
	"github.com/google/uuid"
)

type QueryRepository struct {
	db *database.Queries
}

func NewQueryRepository(db *sql.DB) *QueryRepository {
	return &QueryRepository{
		db: database.New(db),
	}
}

func (repository *QueryRepository) CreateUser(ctx context.Context, createUserRequest *models.CreateUserRequest) (database.User, error) {
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      createUserRequest.User,
		ApiKey:    utils.GenerateApiKey(),
	}

	return repository.db.CreateUser(ctx, params)
}

func (repository *QueryRepository) GetUser(ctx context.Context, apiKey string) (database.User, error) {
	return repository.db.GetUser(ctx, apiKey)
}

func (repository *QueryRepository) CreateFeed(ctx context.Context, user *database.User, request *models.CreateFeedRequest) (database.Feed, error) {
	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      request.Name,
		Url:       request.Url,
		UserID:    user.ID,
	}

	return repository.db.CreateFeed(ctx, createFeedParams)
}

func (repository *QueryRepository) GetFeeds(ctx context.Context) ([]database.Feed, error) {
	return repository.db.GetFeeds(ctx)
}

func (repository *QueryRepository) CreateFeedFollow(ctx context.Context, user *database.User, feedId uuid.UUID) (database.FeedFollow, error) {
	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feedId,
		UserID:    user.ID,
	}

	return repository.db.CreateFeedFollow(ctx, createFeedFollowParams)
}

func (repository *QueryRepository) GetFeedFollows(ctx context.Context, userId uuid.UUID) ([]database.FeedFollow, error) {
	return repository.db.GetFeedFollows(ctx, userId)
}

func (repository *QueryRepository) DeleteFeedFollows(ctx context.Context, feedFollowsId uuid.UUID) error {
	return repository.db.DeleteFeedFollow(ctx, feedFollowsId)
}

func (repository *QueryRepository) GetNextFeedsToFetch(ctx context.Context) ([]database.Feed, error) {
	return repository.db.GetNextFeedsToFetch(ctx)
}

func (repository *QueryRepository) CreatePost(ctx context.Context, request *models.RssItem, feedId uuid.UUID) (database.Post, error) {
	var description sql.NullString
	if request.Description == "" {
		description = sql.NullString{Valid: false}
	} else {
		description = sql.NullString{String: request.Description, Valid: true}
	}

	const layout = "Mon, 02 Jan 2006 15:04:05 -0700"
	pubDate, err := time.Parse(layout, request.PubDate)
	if err != nil {
		return database.Post{}, err
	}

	var publishedAt sql.NullTime
	if request.Description == "" {
		publishedAt = sql.NullTime{Valid: false}
	} else {
		publishedAt = sql.NullTime{Time: pubDate, Valid: true}
	}

	createPostParam := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       request.Title,
		Url:         request.Link,
		Description: description,
		PublishedAt: publishedAt,
		FeedID:      feedId,
	}

	return repository.db.CreatePost(ctx, createPostParam)
}

func (repository *QueryRepository) GetPostsByUser(ctx context.Context, userId uuid.UUID, limit int32) ([]database.Post, error) {
	params := database.GetPostsByUserParams{
		UserID: userId,
		Limit:  limit,
	}

	return repository.db.GetPostsByUser(ctx, params)
}
