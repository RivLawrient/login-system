package auth

import (
	"errors"

	"github.com/RivLawrient/login-system/backend/internal/dto"
	"github.com/RivLawrient/login-system/backend/internal/errs"
	"github.com/RivLawrient/login-system/backend/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Service  *Service
	Validate *validator.Validate
}

func NewHandler(service *Service, validate *validator.Validate) *Handler {
	return &Handler{
		Service:  service,
		Validate: validate,
	}
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	req := dto.RegisterReq{}
	c.ShouldBindJSON(&req)

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(400, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	res, accessToken, refreshToken, err := h.Service.Register(c, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errs.ErrEmailUsed) {
			c.JSON(400, dto.ResponseWeb[map[string]any]{
				Message: "validation failed",
				Data: gin.H{
					"email": err.Error(),
				},
			})
			return
		}

		c.JSON(500, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", *refreshToken, 0, "/", "", false, false)

	c.JSON(200, dto.ResponseWeb[dto.RegisterRes]{
		Message: "success register",
		Data: dto.RegisterRes{
			Email:       res.Email,
			AccessToken: *accessToken,
		},
	})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	req := dto.LoginReq{}
	c.ShouldBindJSON(&req)

	if err := h.Validate.Struct(req); err != nil {
		c.JSON(400, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	res, accessToken, refreshToken, err := h.Service.Login(c, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidEmailPassword) {
			c.JSON(400, dto.ResponseWeb[map[string]any]{
				Message: "validation failed",
				Data: gin.H{
					"email":    err.Error(),
					"password": err.Error(),
				},
			})
			return
		}

		c.JSON(500, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", *refreshToken, 0, "/", "", false, false)

	c.JSON(200, dto.ResponseWeb[dto.RegisterRes]{
		Message: "success login",
		Data: dto.RegisterRes{
			Email:       res.Email,
			AccessToken: *accessToken,
		},
	})
}

func (h *Handler) RefreshHandler(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, dto.ResponseWeb[any]{
			Message: errs.ErrInvalidRefreshToken.Error(),
		})
		return
	}

	accessToken, newRefresh, err := h.Service.RefreshAccessToken(c, token)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidRefreshToken) {
			c.JSON(401, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidRefreshToken.Error(),
			})
			return
		}

		c.JSON(500, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	// set new refresh cookie
	c.SetCookie("refresh_token", *newRefresh, 0, "/", "", false, false)

	c.JSON(200, dto.ResponseWeb[dto.RefreshRes]{
		Message: "success refresh",
		Data: dto.RefreshRes{
			AccessToken: *accessToken,
		},
	})
}

func (h *Handler) MeHandler(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, dto.ResponseWeb[any]{
			Message: errs.ErrInvalidAccessToken.Error(),
		})
		return
	}

	idStr, ok := uid.(string)
	if !ok {
		c.JSON(401, dto.ResponseWeb[any]{
			Message: errs.ErrInvalidAccessToken.Error(),
		})
		return
	}

	user, err := h.Service.Me(c, idStr)
	if err != nil {
		if errors.Is(err, errs.ErrDataNotFound) {
			c.JSON(404, dto.ResponseWeb[any]{
				Message: errs.ErrDataNotFound.Error(),
			})
			return
		}
		c.JSON(500, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.JSON(200, dto.ResponseWeb[dto.MeRes]{
		Message: "success get user",
		Data:    dto.MeRes{ID: user.ID, Email: user.Email},
	})
}
