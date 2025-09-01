package types

import (
	"sync/atomic"

	"github.com/HemahWeb/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	JWTSecret      string
}
