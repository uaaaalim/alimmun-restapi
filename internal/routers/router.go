package routers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/uaaaalim/alimmun-restapi/internal/cache"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, redisClient *redis.Client) *http.ServeMux {
	mux := http.NewServeMux()

	dbServices := *database.NewDatabaseServices(db)
	cacheServices := *cache.NewCacheServices(redisClient)

	_ = dbServices
	_ = cacheServices

	// middlewares

	return mux
}
