package service

import (
	"net/http"

	"github.com/bnakarmi/blog_aggregator/handlers"
	"github.com/bnakarmi/blog_aggregator/repository"
	"github.com/bnakarmi/blog_aggregator/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type RoutingService struct {
	userHandlers *handlers.UserHandler
	feedHandlers *handlers.FeedHandler
}

func NewRoutingService(repository *repository.QueryRepository) *RoutingService {
	userHandlers := &handlers.UserHandler{Repository: repository}
	feedHandlers := &handlers.FeedHandler{Repository: repository}

	return &RoutingService{
		userHandlers: userHandlers,
		feedHandlers: feedHandlers,
	}
}

func (rs *RoutingService) Initialize() *chi.Mux {
	userHandlers := rs.userHandlers
	feedHandlers := rs.feedHandlers

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Post("/v1/user", userHandlers.CreateUserHandler)
	router.Get("/v1/user", authMiddleware(userHandlers.GetUserHandler))

	router.Post("/v1/feeds", authMiddleware(feedHandlers.CreateFeedHandler))
	router.Get("/v1/feeds", feedHandlers.GetFeedsHandler)

	router.Post("/v1/feeds/follows", authMiddleware(feedHandlers.CreateFeedFollowsHandler))
	router.Get("/v1/feeds/follows", authMiddleware(feedHandlers.GetFeedFollowsHandler))
	router.Delete("/v1/feeds/follows/{feedFollowsId}", authMiddleware(feedHandlers.DeleteFeedFollowsHandler))

	router.Get("/v1/posts", authMiddleware(userHandlers.GetPostsByUser))

	return router
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.GetApiKey(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
