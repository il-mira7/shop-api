package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageService interface {
	Create(ctx context.Context, image []byte, productID uuid.UUID) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, image []byte) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]byte, error)
	GetByImageID(ctx context.Context, id uuid.UUID) ([]byte, error)
}

type imageHandler struct {
	service ImageService
}

func NewImageHandler(service ImageService) *imageHandler {
	return &imageHandler{service: service}
}

// Create godoc
// @Summary      добавление изображения
// @Description  на вход подается byte array изображения и id товара
// @Tags         images
// @Accept       octet-stream
// @Produce      json
// @Param        product_id  path      string  true  "ID товара"
// @Param        image       body      string  true  "Бинарный файл"
// @Success      201         {string}  string  "ID изображения"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /images/product/{product_id} [post]
func (h *imageHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	productID, ok := parseIDFromParam(c, "product_id")
	if !ok {
		return
	}

	image, err := c.GetRawData()
	if err != nil {
		sendError(c, err)
		return
	}

	id, err := h.service.Create(ctx, image, productID)
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, id)
}

// Update godoc
// @Summary      Изменение изображения
// @Description  на вход подается id изображения и новая строка для замены
// @Tags         images
// @Accept       octet-stream
// @Produce      json
// @Param        id     path  string  true  "ID изображения"
// @Param        image  body  string  true  "Новые бинарные данные изображения"
// @Success      200    "Успешное обновление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /images/{id} [put]
func (h *imageHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	image, err := c.GetRawData()
	if err != nil {
		sendError(c, err)
		return
	}

	err = h.service.Update(ctx, id, image)
	if err != nil {
		sendError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary      Удалить изображение
// @Description  Удаление изображения по id изображения
// @Tags         images
// @Param        id   path  string  true  "ID изображения (UUID)"
// @Success      204  "Успешное удаление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /images/{id} [delete]
func (h *imageHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		sendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByProductID godoc
// @Summary      Получение изображения по id товара
// @Description  Получение изображения конкретного товара по id товара
// @Tags         images
// @Produce      octet-stream
// @Param        product_id  path  string  true  "ID товара"
// @Success      200         {file}  binary
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /images/product/{product_id} [get]
func (h *imageHandler) GetByProductID(c *gin.Context) {
	ctx := c.Request.Context()

	productID, ok := parseIDFromParam(c, "product_id")
	if !ok {
		return
	}

	image, err := h.service.GetByProductID(ctx, productID)
	if err != nil {
		sendError(c, err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"product_%s.jpg\"", productID.String()))
	c.Data(http.StatusOK, "application/octet-stream", image)
}

// GetByImageID godoc
// @Summary      Получение изображения по id изображения
// @Description  Получение изображения по id изображения
// @Tags         images
// @Produce      octet-stream
// @Param        id   path  string  true  "ID изображения"
// @Success      200  {file}  binary
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /images/{id} [get]
func (h *imageHandler) GetByImageID(c *gin.Context) {
	ctx := c.Request.Context()

	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	image, err := h.service.GetByImageID(ctx, id)
	if err != nil {
		sendError(c, err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"image_%s.jpg\"", id.String()))
	c.Data(http.StatusOK, "application/octet-stream", image)
}
