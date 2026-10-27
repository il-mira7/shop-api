package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserClient interface {
	Register(ctx context.Context, email, name, surname, phone, password string) (string, error)
	Authenticate(ctx context.Context, email, password string) (string, error)
	ResetPassword(ctx context.Context, email string) (bool, error)
}

type userHandler struct {
	client UserClient
}

func NewUserHandler(client UserClient) *userHandler {
	return &userHandler{client: client}
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ResetResponse struct {
	Success bool `json:"success"`
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт учётную запись в сервисе авторизации и возвращает JWT-токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Данные для регистрации"
// @Success      201 {object} TokenResponse
// @Failure      400 {object} errorResponse "Некорректные данные"
// @Failure      409 {object} errorResponse "Email уже зарегистрирован"
// @Router       /register [post]
func (h *userHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	token, err := h.client.Register(c.Request.Context(), req.Email, req.Name, req.Surname, req.PhoneNumber, req.Password)
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{Token: token})
}

// Authenticate godoc
// @Summary      Вход
// @Description  Проверяет логин и пароль, возвращает JWT-токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body AuthRequest true "Логин и пароль"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} errorResponse "Некорректные данные"
// @Failure      401 {object} errorResponse "Неверный логин или пароль"
// @Router       /auth [post]
func (h *userHandler) Authenticate(c *gin.Context) {
	var req AuthRequest
	if !bindJSON(c, &req) {
		return
	}

	token, err := h.client.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: token})
}

// ResetPassword godoc
// @Summary      Восстановление пароля
// @Description  Генерирует новый пароль и "отправляет на почту" (заглушка: вывод в консоль auth-сервиса)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ResetRequest true "Email пользователя"
// @Success      200 {object} ResetResponse
// @Failure      400 {object} errorResponse "Некорректные данные"
// @Router       /reset [post]
func (h *userHandler) ResetPassword(c *gin.Context) {
	var req ResetRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.client.ResetPassword(c.Request.Context(), req.Email)
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, ResetResponse{Success: res})
}
