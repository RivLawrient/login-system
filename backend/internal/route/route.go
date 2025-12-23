package route

import (
	"net/http"

	"github.com/RivLawrient/login-system/backend/internal/dto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	App *gin.Engine
}

func (r Routes) SetupRouts() {
	r.App.Use(cors.Default())

	r.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.ResponseWeb[any]{
			Message: "This your root api",
		})
	})
}
