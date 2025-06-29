package handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rohitsh16/go-blog/backend/server/config"
)

type Service struct {
	Config       *config.Config
	Mysql        *sql.DB
	GinFramework *gin.Engine
}
