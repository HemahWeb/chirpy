package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/HemahWeb/chirpy/internal/auth"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responseVals struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := h.config.DB.GetUserByEmailForAuth(r.Context(), params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, h.config.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	refreshToken, err := h.config.DB.CreateRefreshToken(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, responseVals{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
