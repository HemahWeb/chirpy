package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HemahWeb/chirpy/internal/utils"
	"github.com/google/uuid"
)

func (h *Handler) PolkaUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var params parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		utils.RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	_, err = h.config.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	if params.Event == "user.upgraded" {
		err = h.config.DB.UpgradeUserToChirpyRed(r.Context(), userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user", err)
			return
		}
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}
