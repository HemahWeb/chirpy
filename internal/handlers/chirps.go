package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/HemahWeb/chirpy/internal/auth"
	"github.com/HemahWeb/chirpy/internal/database"
	"github.com/HemahWeb/chirpy/internal/types"
	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) PostChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, h.config.JWTSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	chirp, err := h.config.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp: "+err.Error(), err)
		return
	}

	// map DB -> API (stable keys, decoupled from schema)
	utils.RespondWithJSON(w, http.StatusCreated, types.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (h *Handler) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := h.config.DB.GetChirps(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	// map DB -> API (stable keys, decoupled from schema)
	chirpsAPI := []types.Chirp{}
	for _, chirp := range chirps {
		chirpsAPI = append(chirpsAPI, types.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, chirpsAPI)
}

func (h *Handler) GetChirpsByID(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID: "+err.Error(), err)
		return
	}

	chirp, err := h.config.DB.GetChirpByID(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Chirp with ID "+id.String()+" does not exist", err)
		return
	}

	// map DB -> API (stable keys, decoupled from schema)
	utils.RespondWithJSON(w, http.StatusOK, types.Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
