package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HemahWeb/chirpy/internal/auth"
	"github.com/HemahWeb/chirpy/internal/database"
	"github.com/HemahWeb/chirpy/internal/types"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) UsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// SQLC returns a database.User without API JSON tags
	user, err := h.config.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	// Map DB -> API (stable keys, decoupled from schema)
	utils.RespondWithJSON(w, http.StatusCreated, types.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (h *Handler) UsersReset(w http.ResponseWriter, r *http.Request) {
	if h.config.Platform == "" {
		utils.RespondWithError(w, http.StatusServiceUnavailable, "Platform not set in config", nil)
		return
	}
	if h.config.Platform != "dev" {
		utils.RespondWithError(w, http.StatusForbidden, "Forbidden in non-dev platform", nil)
		return
	}

	err := h.config.DB.ResetUsers(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to reset users in database:", err)
		return
	}

	h.config.FileserverHits.Store(0)

	utils.RespondWithJSON(w, http.StatusOK, nil)
}

func (h *Handler) UsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't get token: "+err.Error(), err)
		return
	}

	userID, err := auth.ValidateJWT(token, h.config.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	var params parameters
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := h.config.DB.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, types.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
