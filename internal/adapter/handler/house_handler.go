package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type HouseHandler struct {
	houseUseCase port.HouseUseCase
}

func NewHouseHandler(e *echo.Group, houseUseCase port.HouseUseCase) {
	handler := &HouseHandler{houseUseCase: houseUseCase}
	e.POST("/houses", handler.Create)
	e.GET("/houses/:id", handler.GetByID)
	e.GET("/houses", handler.List)
	e.PUT("/houses/:id", handler.Update)

	// Chỉ người có quyền house:delete mới được xoá nhà
	e.DELETE("/houses/:id", handler.Delete, PermissionMiddleware("house:delete"))
}

func (h *HouseHandler) Create(c echo.Context) error {
	house := new(domain.House)
	if err := c.Bind(house); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.houseUseCase.CreateHouse(c.Request().Context(), house)
	if err != nil {
		log.Printf("Error creating house: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, house)
}

func (h *HouseHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}

	house, err := h.houseUseCase.GetHouse(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "house not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, house)
}

func (h *HouseHandler) List(c echo.Context) error {
	// Lấy tham số phân trang từ URL, ví dụ: /api/v1/houses?cursor=15&limit=10
	cursorStr := c.QueryParam("cursor")
	limitStr := c.QueryParam("limit")

	cursor := 0
	limit := 10

	if cursorStr != "" {
		if parsed, err := strconv.Atoi(cursorStr); err == nil {
			cursor = parsed
		}
	}
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	houses, err := h.houseUseCase.ListHouses(c.Request().Context(), cursor, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	// Xác định cursor cho trang tiếp theo (là ID của item cuối cùng trong danh sách)
	var nextCursor int64 = 0
	if len(houses) > 0 {
		nextCursor = houses[len(houses)-1].ID
	}

	// Trả về dữ liệu bọc trong một struct để dễ tích hợp với Vue.js
	response := map[string]interface{}{
		"data":        houses,
		"next_cursor": nextCursor,
		"has_more":    len(houses) == limit, // Nếu số lượng trả về bằng limit, khả năng cao là còn trang sau
	}

	return c.JSON(http.StatusOK, response)
}

func (h *HouseHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}

	house := new(domain.House)
	if err := c.Bind(house); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	house.ID = id

	err = h.houseUseCase.UpdateHouse(c.Request().Context(), house)
	if err != nil {
		if err.Error() == "house not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "house updated successfully"})
}

func (h *HouseHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}

	err = h.houseUseCase.DeleteHouse(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "house not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "house deleted successfully"})
}
