package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/HemahWeb/chirpy/internal/auth"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Refresh token invalid: "+err.Error(), err)
		return
	}

	user, err := h.config.DB.GetUserIDFromRefreshToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusUnauthorized, "Refresh token not in database", err)
			return
		}
		utils.RespondWithError(w, http.StatusUnauthorized, "Could not get user from refresh token", err)
		return
	}

	if user.RevokedAt.Valid {
		utils.RespondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	if user.ExpiresAt.Before(time.Now()) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	// Generate a new JWT token for the user and return it in the response
	tokenString, err := auth.MakeJWT(user.UserID, h.config.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	type responseVals struct {
		Token string `json:"token"`
	}
	utils.RespondWithJSON(w, http.StatusOK, responseVals{
		Token: tokenString,
	})
}
