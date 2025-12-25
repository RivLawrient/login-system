package middleware

import (
	"net/http"
	"strings"

	"github.com/RivLawrient/login-system/backend/internal/dto"
	"github.com/RivLawrient/login-system/backend/internal/errs"
	"github.com/RivLawrient/login-system/backend/internal/helper"
	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := helper.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			ctx.Abort()
			return
		}

		idVal, ok := claims["user_id"].(string)
		if !ok || idVal == "" {
			ctx.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", idVal)
		ctx.Next()
	}
}
