package handler

import (
	"context"
	"net/http"
	"shopAPI/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductService interface {
	Create(ctx context.Context, product *domain.Product) error
	DecreaseStock(ctx context.Context, id uuid.UUID, count int) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type productHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *productHandler {
	return &productHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Добавление товара
// @Description  Добавление товара на вход подается json.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  	body   		productDTO    true 	"Данные нового товара"
// @Success      201   		{object}  	productDTO       	"Возвращает созданный товар"
// @Failure      400  		{object}  	errorResponse  		"Валидация прошла неуспешно"
// @Failure      500  		{object}  	errorResponse  		"Внутренняя ошибка сервера"
// @Router       /products [post]
func (h *productHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var dto productDTO
	if !bindJSON(c, &dto) {
		return
	}

	product := dto.toDomain()
	if err := h.service.Create(ctx, product); err != nil {
		sendError(c, err)
		return
	}

	dto.fromDomain(product)
	c.JSON(http.StatusCreated, dto)
}

// DecreaseStock godoc
// @Summary      Уменьшение количества товара
// @Description  Уменьшение количества товара (на вход запросу подается id товара и на сколько уменьшить).
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id      path      string      			 true 	"ID товара" format(uuid)
// @Param        input 	 body      decreaseStockRequest  true  	"Количество, на сколько уменьшить"
// @Success      200  "Успешное обновление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно или недостаточно товара на складе"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /products/{id}/stock [patch]
func (h *productHandler) DecreaseStock(c *gin.Context) {
	ctx := c.Request.Context()

	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	var countReq decreaseStockRequest
	if !bindJSON(c, &countReq) {
		return
	}

	if err := h.service.DecreaseStock(ctx, id, countReq.Count); err != nil {
		sendError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// GetByID godoc
// @Summary      Получение товара по id
// @Description  Получение товара по его идентификатору
// @Tags         products
// @Produce      json
// @Param        id      path      string      true 	"ID товара" format(uuid)
// @Success      200  {object}  productDTO     "Товар с указанным id"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /products/{id} [get]
func (h *productHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	product, err := h.service.GetByID(ctx, id)
	if err != nil {
		sendError(c, err)
		return
	}

	var dto productDTO
	dto.fromDomain(product)

	c.JSON(http.StatusOK, dto)
}

// GetAll godoc
// @Summary      Получение всех доступных товаров
// @Description  Получение всех товаров, предусмотрены опциональные параметры пагинации в строке запроса: limit и offset. В случае отсутствия этих параметров возвращается весь список.
// @Tags         products
// @Produce      json
// @Param        limit   query     int  false  "Количество записей"
// @Param        offset  query     int  false  "Смещение"
// @Success      200  	{object}  productsListResponse
// @Failure      400  	{object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      500 	{object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /products [get]
func (h *productHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	limit, offset, ok := parsePagination(c)
	if !ok {
		return
	}

	products, count, err := h.service.GetAll(ctx, limit, offset)
	if err != nil {
		sendError(c, err)
		return
	}

	dto := make([]productDTO, len(products))
	for i, v := range products {
		dto[i].fromDomain(v)
	}

	c.JSON(http.StatusOK, productsListResponse{
		Products: dto,
		Total:    count,
		Limit:    limit,
		Offset:   offset,
	})
}

// Delete godoc
// @Summary      Удаление товара
// @Description  Удаление товара по его идентификатору
// @Tags         products
// @Param        id   path  string  true  "ID товара" format(uuid)
// @Success      204  "Успешное удаление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /products/{id} [delete]
func (h *productHandler) Delete(c *gin.Context) {
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
