package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/eral-promo-library-go/library/data"
	"github.com/fritz-immanuel/eral-promo-library-go/src/app/web"
	"github.com/jmoiron/sqlx"
)

// RegisterWebRoutes  is a function to register all WEB Routes in the projectbase
func RegisterWebRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine) {
	v1 := router.Group("/web/v1")
	{
		web.RegisterRoutes(db, dataManager, router, v1)
	}
}
