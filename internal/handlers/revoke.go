package handlers

import (
	"net/http"

	"github.com/HemahWeb/chirpy/internal/auth"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't get token: "+err.Error(), err)
		return
	}

	err = h.config.DB.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't revoke token", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
