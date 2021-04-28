package middlewaresv1

import (
	"net/http"

	"github.com/fittipaldi/go-r3-d3/config"

	"github.com/gorilla/context"
	"gorm.io/gorm"
)

type CommonMiddlewareV3 struct {
	GormDB *gorm.DB       `json:"gorm_db"`
	Config *config.Config `json:"config"`
}

// Middleware function, which will be called for each request
func (mw *CommonMiddlewareV3) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "gormDB", mw.GormDB)
		context.Set(r, "config", mw.Config)
		next.ServeHTTP(w, r)
	})
}
