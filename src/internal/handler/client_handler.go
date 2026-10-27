package handler

import (
	"context"
	"net/http"
	"shopAPI/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientService interface {
	Create(ctx context.Context, client *domain.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByNames(ctx context.Context, firstName string, lastName string) ([]*domain.Client, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*domain.Client, int, error)
	Update(ctx context.Context, id uuid.UUID, newAddress *domain.Address) (*domain.Client, error)
}

type clientHandler struct {
	service ClientService
}

func NewClientHandler(service ClientService) *clientHandler {
	return &clientHandler{service: service}
}

// Create godoc
// @Summary      Добавление клиента
// @Description  Добавление клиента на вход подается json.
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        client  body     clientDTO  true  	"Данные нового клиента"
// @Success      201    {object}  clientDTO         "Возвращает созданного клиента"
// @Failure      400  	{object}  errorResponse  	"Валидация прошла неуспешно"
// @Failure      500  	{object}  errorResponse  	"Внутренняя ошибка сервера"
// @Router       /clients [post]
func (h *clientHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var dto clientDTO
	if !bindJSON(c, &dto) {
		return
	}

	client := dto.toDomain()

	if err := h.service.Create(ctx, client); err != nil {
		sendError(c, err)
		return
	}

	dto.fromDomain(client)

	c.JSON(http.StatusCreated, dto)
}

// Delete godoc
// @Summary      Удаление клиента
// @Description  Удаление клиента по его идентификатору.
// @Tags         clients
// @Param        id   path  string  true  "ID клиента" format(uuid)
// @Success      204  "Успешное удаление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /clients/{id} [delete]
func (h *clientHandler) Delete(c *gin.Context) {
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

// GetByNames godoc
// @Summary      Получение клиентов по имени и фамилии
// @Description  Ищет клиентов по точному совпадению имени и фамилии. Если ничего не найдено, возвращается пустой список.
// @Tags         clients
// @Produce      json
// @Param        name     query     string  true  "Имя клиента"
// @Param        surname  query     string  true  "Фамилия клиента"
// @Success      200     {object}  clientsListResponse
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /clients/search [get]
func (h *clientHandler) GetByNames(c *gin.Context) {
	ctx := c.Request.Context()
	firstName := c.Query("name")
	lastName := c.Query("surname")
	if firstName == "" || lastName == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "query parameters name and surname are required"})
		return
	}

	clients, err := h.service.GetByNames(ctx, firstName, lastName)
	if err != nil {
		sendError(c, err)
		return
	}
	dto := make([]clientDTO, len(clients))
	for i, v := range clients {
		dto[i].fromDomain(v)
	}

	c.JSON(http.StatusOK, clientsListResponse{
		Clients: dto,
		Total:   len(dto),
		Limit:   len(dto),
		Offset:  0,
	})
}

// GetAll godoc
// @Summary      Получение всех клиентов
// @Description  Получение всех клиентов, предусмотрены опциональные параметры пагинации в строке запроса: limit и offset. В случае отсутствия этих параметров возвращается весь список.
// @Tags         clients
// @Produce      json
// @Param        limit   query     int  false  "Количество записей"
// @Param        offset  query     int  false  "Смещение"
// @Success      200     {object}  clientsListResponse
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /clients [get]
func (h *clientHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	limit, offset, ok := parsePagination(c)
	if !ok {
		return
	}

	clients, count, err := h.service.GetAll(ctx, limit, offset)
	if err != nil {
		sendError(c, err)
		return
	}

	dto := make([]clientDTO, len(clients))
	for i, v := range clients {
		dto[i].fromDomain(v)
	}

	c.JSON(http.StatusOK, clientsListResponse{
		Clients: dto,
		Total:   count,
		Limit:   limit,
		Offset:  offset,
	})
}

// Update godoc
// @Summary      Изменение адреса
// @Description  Изменение адреса клиента (параметры: id и новый адрес в виде json).
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id      path      string      true "ID клиента" format(uuid)
// @Param        address body      addressDTO  true  "Новый адрес клиента"
// @Success      200  {object}  clientDTO      "Успешное обновление"
// @Failure      400  {object}  errorResponse  "Валидация прошла неуспешно"
// @Failure      404  {object}  errorResponse  "Данные отсутствуют"
// @Failure      500  {object}  errorResponse  "Внутренняя ошибка сервера"
// @Router       /clients/{id}/address [patch]
func (h *clientHandler) Update(c *gin.Context) {
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

	client, err := h.service.Update(ctx, id, address)
	if err != nil {
		sendError(c, err)
		return
	}

	var cDTO clientDTO
	cDTO.fromDomain(client)

	c.JSON(http.StatusOK, cDTO)
}
