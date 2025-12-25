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

	res, err := h.Service.Register(c, req.Email, req.Password)
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

		c.JSON(400, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.JSON(200, dto.ResponseWeb[dto.RegisterRes]{
		Message: "success",
		Data: dto.RegisterRes{
			Email: res.Email,
		},
	})
}
