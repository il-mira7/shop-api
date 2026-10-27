package handler

import (
	"context"
	"net/http"
	"shopAPI/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SupplierService interface {
	Create(ctx context.Context, supplier *domain.Supplier) error
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Supplier, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) (*domain.Supplier, error)
}

type supplierHandler struct {
	service SupplierService
}

func NewSupplierHandler(service SupplierService) *supplierHandler {
	return &supplierHandler{service: service}
}

// Create godoc
// @Summary      Добавление поставщика
// @Description  Добавление поставщика на вход подается json.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        supplier  body   supplierDTO  true "Данные нового поставщика"
// @Success      201    {object}  supplierDTO       "Возвращает созданного поставщика"
// @Failure      400  	{object}  errorResponse  	"Валидация прошла неуспешно"
// @Failure      500  	{object}  errorResponse  	"Внутренняя ошибка сервера"
// @Router       /suppliers [post]
func (h *supplierHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var dto supplierDTO
	if !bindJSON(c, &dto) {
		return
	}

	supplier := dto.toDomain()
	if err := h.service.Create(ctx, supplier); err != nil {
		sendError(c, err)
		return
	}

	dto.fromDomain(supplier)

	c.JSON(http.StatusCreated, dto)

}

// Update godoc
// @Summary      Изменение адреса
// @Description  Изменение адреса поставщика (параметры: id и новый адрес в виде json).
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        id      path      string      true "ID поставщика" format(uuid)
// @Param        address body      addressDTO  true  "Новый адрес поставщика"
// @Success      200  {object}  supplierDTO    "Успешное обновление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /suppliers/{id}/address [patch]
func (h *supplierHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	var dto addressDTO
	if !bindJSON(c, &dto) {
		return
	}

	address := dto.toDomain()
	supplier, err := h.service.Update(ctx, id, address)
	if err != nil {
		sendError(c, err)
		return
	}

	var supplierDto supplierDTO
	supplierDto.fromDomain(supplier)

	c.JSON(http.StatusOK, supplierDto)

}

// GetByID godoc
// @Summary      Получение поставщика по id
// @Description  Получение поставщика по его идентификатору
// @Tags         suppliers
// @Produce      json
// @Param        id      path      string      true "ID поставщика" format(uuid)
// @Success      200  {object}  supplierDTO    "Поставщик с указанным id"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /suppliers/{id} [get]
func (h *supplierHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, ok := parseIDFromParam(c, "id")
	if !ok {
		return
	}

	supplier, err := h.service.GetByID(ctx, id)
	if err != nil {
		sendError(c, err)
		return
	}

	var dto supplierDTO

	dto.fromDomain(supplier)

	c.JSON(http.StatusOK, dto)
}

// GetAll godoc
// @Summary      Получение всех поставщиков
// @Description  Получение всех поставщиков, предусмотрены опциональные параметры пагинации в строке запроса: limit и offset. В случае отсутствия этих параметров возвращается весь список.
// @Tags         suppliers
// @Produce      json
// @Param        limit   query     int  false  "Количество записей"
// @Param        offset  query     int  false  "Смещение"
// @Success      200  	{object}  suppliersListResponse
// @Failure      400  	{object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      500 	{object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /suppliers [get]
func (h *supplierHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	limit, offset, ok := parsePagination(c)
	if !ok {
		return
	}

	suppliers, count, err := h.service.GetAll(ctx, limit, offset)
	if err != nil {
		sendError(c, err)
		return
	}

	dto := make([]supplierDTO, 0, len(suppliers))
	for _, s := range suppliers {
		var d supplierDTO
		d.fromDomain(s)
		dto = append(dto, d)
	}
	c.JSON(http.StatusOK, suppliersListResponse{
		Suppliers: dto,
		Total:     count,
		Limit:     limit,
		Offset:    offset,
	})
}

// Delete godoc
// @Summary      Удаление поставщика
// @Description  Удаление поставщика по его идентификатору.
// @Tags         suppliers
// @Param        id   path  string  true  "ID поставщика" format(uuid)
// @Success      204  "Успешное удаление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /suppliers/{id} [delete]
func (h *supplierHandler) Delete(c *gin.Context) {
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
