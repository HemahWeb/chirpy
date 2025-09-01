package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HemahWeb/chirpy/internal/auth"
	"github.com/HemahWeb/chirpy/internal/types"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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

	// Respond with user, without hashed password
	utils.RespondWithJSON(w, http.StatusOK, types.User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
