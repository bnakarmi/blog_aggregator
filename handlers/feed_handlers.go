package handlers

import (
	"log"
	"net/http"

	"github.com/bnakarmi/blog_aggregator/internal/models"
	"github.com/bnakarmi/blog_aggregator/repository"
	"github.com/bnakarmi/blog_aggregator/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FeedHandler struct {
	Repository *repository.QueryRepository
}

const CreateFeedErrorMsg = "Unable to create feed"
const CreateFeedFollowErrorMsg = "Error creating feed follows"
const GetFeedsErrorMsg = "Unable to get feeds"
const GetFeedFollowsErrorMsg = "Error getting feed follows"
const InvalidFeedFollowIdMsg = "Invalid feed follows id"
const FeedFollowsDeleteFailedMsg = "Error deleting feed follows"

func (h *FeedHandler) CreateFeedHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, _ := utils.GetApiKey(r)
	user, err := h.Repository.GetUser(r.Context(), apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	var createFeedRequest models.CreateFeedRequest
	err = utils.DecodeJSONFromBody(w, r, &createFeedRequest)
	if err != nil {
		return
	}

	createdFeed, err := h.Repository.CreateFeed(r.Context(), &user, &createFeedRequest)
	if err != nil {
		log.Printf("%s %v", CreateFeedErrorMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, CreateFeedErrorMsg)
		return
	}

	createFeedFollows, err := h.Repository.CreateFeedFollow(r.Context(), &user, createdFeed.ID)
	if err != nil {
		log.Printf("%s %v", CreateFeedFollowErrorMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, CreateFeedFollowErrorMsg)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.CreateFeedResponse{
		Feed:       utils.GetFeedResponse(&createdFeed),
		FeedFollow: createFeedFollows,
	})
}

func (h *FeedHandler) GetFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.Repository.GetFeeds(r.Context())
	if err != nil {
		log.Printf("%s %v", GetFeedsErrorMsg, err)
		utils.RespondWithError(w, http.StatusNotFound, GetFeedsErrorMsg)
		return
	}

    var feedsJson = make([]models.Feed, 0)
    for _, feed := range feeds {
        feedsJson = append(feedsJson, utils.GetFeedResponse(&feed))
    }

	utils.RespondWithJSON(w, http.StatusOK, feedsJson)
}

func (h *FeedHandler) CreateFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, _ := utils.GetApiKey(r)
	user, err := h.Repository.GetUser(r.Context(), apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	var createFeedFollowsRequest models.CreateFeedFollowsRequest
	err = utils.DecodeJSONFromBody(w, r, &createFeedFollowsRequest)
	if err != nil {
		return
	}

	createFeedFollowsResponse, err := h.Repository.CreateFeedFollow(r.Context(), &user, createFeedFollowsRequest.FeedID)
	if err != nil {
		log.Printf("%s %v", CreateFeedFollowErrorMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, CreateFeedFollowErrorMsg)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, createFeedFollowsResponse)
}

func (h *FeedHandler) GetFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, _ := utils.GetApiKey(r)
	user, err := h.Repository.GetUser(r.Context(), apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	feedFollowsList, err := h.Repository.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		log.Printf("%s %v", GetFeedFollowsErrorMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, GetFeedFollowsErrorMsg)
		return
	}

	if feedFollowsList == nil {
		utils.RespondWithJSON(w, http.StatusOK, make([]int, 0))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, feedFollowsList)
}

func (h *FeedHandler) DeleteFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	feedFollowsId, err := uuid.Parse(chi.URLParam(r, "feedFollowsId"))
	if err != nil {
		log.Printf("%s %v", InvalidFeedFollowIdMsg, err)
		utils.RespondWithError(w, http.StatusBadRequest, InvalidFeedFollowIdMsg)
		return
	}

	err = h.Repository.DeleteFeedFollows(r.Context(), feedFollowsId)
	if err != nil {
		log.Printf("%s %v", FeedFollowsDeleteFailedMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, FeedFollowsDeleteFailedMsg)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Feed follows deleted successfully",
	})
}
