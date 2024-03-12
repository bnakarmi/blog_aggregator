package handlers

import (
	"log"
	"net/http"

	"github.com/bnakarmi/blog_aggregator/internal/models"
	"github.com/bnakarmi/blog_aggregator/repository"
	"github.com/bnakarmi/blog_aggregator/utils"
)

type UserHandler struct {
	Repository *repository.QueryRepository
}

const CreateUserFailedMsg = "Create user failed"
const GetPostsFailedMsg = "Error getting user posts"

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest models.CreateUserRequest

	err := utils.DecodeJSONFromBody(w, r, &createUserRequest)
	if err != nil {
		return
	}

	newUser, err := h.Repository.CreateUser(r.Context(), &createUserRequest)
	if err != nil {
		log.Printf("%s %v", CreateUserFailedMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, CreateUserFailedMsg)

		return
	}

	utils.RespondWithJSON(w, http.StatusOK, newUser)
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, _ := utils.GetApiKey(r)
	user, err := h.Repository.GetUser(r.Context(), apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	apiKey, _ := utils.GetApiKey(r)
	user, err := h.Repository.GetUser(r.Context(), apiKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	posts, err := h.Repository.GetPostsByUser(r.Context(), user.ID, 5)
	if err != nil {
		log.Printf("%s %v", GetPostsFailedMsg, err)
		utils.RespondWithError(w, http.StatusInternalServerError, GetPostsFailedMsg)

		return
	}

	postsJson := make([]models.PostResponse, 0)
	for _, post := range posts {
		postsJson = append(postsJson, utils.GetPostResponse(&post))
	}

	utils.RespondWithJSON(w, http.StatusOK, postsJson)
}
