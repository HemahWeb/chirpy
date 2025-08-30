package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HemahWeb/chirpy/internal/utils"
)

func (h *Handler) ChirpValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := utils.ValidateChirp(params.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't validate chirp: "+err.Error(), nil)
		return
	}

	cleanedBody := returnVals{
		Body: cleaned,
	}
	utils.RespondWithJSON(w, http.StatusOK, cleanedBody)
}
