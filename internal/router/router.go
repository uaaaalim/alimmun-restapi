package router

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/uaaaalim/alimmun-restapi/internal/cache"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"github.com/uaaaalim/alimmun-restapi/internal/handlers/auth"
	"github.com/uaaaalim/alimmun-restapi/internal/middlewares"
	"gorm.io/gorm"
)

func NewGeneralRouter(db *gorm.DB, redisClient *redis.Client) *http.Handler {
	mux := http.NewServeMux()

	// services
	dbService := database.NewDatabaseServices(db)
	cacheService := cache.NewCacheServices(redisClient)

	// routes
	auth.NewHandler(dbService, cacheService).RegisterRoutes(mux)

	// middlewares
	handler := middlewares.Logger(mux)

	return &handler
}
