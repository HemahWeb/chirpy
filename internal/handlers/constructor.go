package handlers

import (
	"github.com/HemahWeb/chirpy/internal/types"
)

type Handler struct {
	config *types.ApiConfig
}

func New(config *types.ApiConfig) *Handler {
	return &Handler{config: config}
}
