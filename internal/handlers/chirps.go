package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/HemahWeb/chirpy/internal/database"
	"github.com/HemahWeb/chirpy/internal/types"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) ChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// validate chirp for length and bad words
	cleaned, err := utils.ValidateChirp(params.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't validate chirp: "+err.Error(), nil)
		return
	}

	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID: "+err.Error(), err)
		return
	}

	chirp, err := h.config.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp: "+err.Error(), err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, types.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
