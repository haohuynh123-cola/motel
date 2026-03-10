package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type UserHandler struct {
	userUseCase port.UserUseCase
}

func NewUserHandler(e *echo.Group, uc port.UserUseCase) {
	handler := &UserHandler{userUseCase: uc}
	e.POST("/auth/register", handler.Register)
	e.POST("/auth/login", handler.Login)
}

func NewProtectedUserHandler(e *echo.Group, uc port.UserUseCase) {
	handler := &UserHandler{userUseCase: uc}
	e.GET("/auth/me", handler.GetMe)
}

func (h *UserHandler) Register(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.userUseCase.Register(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user.Password = ""
	return c.JSON(http.StatusCreated, user)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.userUseCase.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) GetMe(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*jwt.MapClaims)

	idFloat := (*claims)["id"].(float64)
	userID := int64(idFloat)

	user, err := h.userUseCase.GetUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
