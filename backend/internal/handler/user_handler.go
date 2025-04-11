package handler

import (
	"net/http"

	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUsecase    *usecase.IUserUsecase
	UserRepository repository.IUserRepository
}

// TODO: usecaseを受け取るようにしてより疎結合にする
func NewUserHandler(
	repo repository.IUserRepository,
	hasher service.Hasher,
	tokenService service.TokenService,
) UserHandler {
	return UserHandler{
		UserUsecase: usecase.NewUserUsecase(
			repo,
			hasher,
			tokenService,
		),
		UserRepository: repo,
	}
}

func (h *UserHandler) Register(g *echo.Group) {
	g.POST("/signup", h.SignUp)
	g.POST("/login", h.Login)
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignUp(c echo.Context) error {
	var params SignUpRequest

	if err := c.Bind(&params); err != nil {
		return c.JSON(400, map[string]any{
			"error": err.Error(),
		})
	}

	user, err := h.UserUsecase.SignUp(params.Name, params.Email, params.Password)
	if err != nil {
		return c.JSON(400, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(200, map[string]any{
		"user": user,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	User  LoginUserResponse `json:"user"`
	Token string            `json:"token"`
}

func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]any{
			"error": err.Error(),
		})
	}

	token, err := h.UserUsecase.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return c.JSON(400, map[string]any{
			"error": err.Error(),
		})
	}

	user, err := h.UserRepository.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(400, map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		User: LoginUserResponse{
			ID:    user.GetID(),
			Name:  user.GetName(),
			Email: user.GetEmail(),
		},
		Token: token,
	})
}
