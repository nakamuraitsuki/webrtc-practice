package handler

import (
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUsecase *usecase.IUserUsecase
}

func NewUserHandler(repo repository.IUserRepository, hasher service.Hasher, tokenService service.TokenService) UserHandler {
	return UserHandler{
		UserUsecase: usecase.NewUserUsecase(
			repo,
			hasher,
			tokenService,
		),
	}
}

func (h *UserHandler) Register(g *echo.Group) {
	g.POST("/register", h.RegisterUser)
	g.POST("/authenticate", h.AuthenticateUser)
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.UserUsecase.RegisterUser(name, email, password)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(200, map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandler) AuthenticateUser(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	token, err := h.UserUsecase.AuthenticateUser(email, password)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(200, map[string]interface{}{
		"token": token,
	})
}
