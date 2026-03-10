package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

type RoomHandler struct {
	roomUseCase port.RoomUseCase
}

func NewRoomHandler(e *echo.Group, uc port.RoomUseCase) {
	handler := &RoomHandler{roomUseCase: uc}
	e.POST("/rooms", handler.Create)
	e.GET("/rooms/:id", handler.GetByID)
	e.GET("/houses/:house_id/rooms", handler.ListByHouseID)
	e.PUT("/rooms/:id", handler.Update)
	e.DELETE("/rooms/:id", handler.Delete)
}

func (h *RoomHandler) Create(c echo.Context) error {
	room := new(domain.Room)
	if err := c.Bind(room); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.roomUseCase.CreateRoom(c.Request().Context(), room)
	if err != nil {
		log.Printf("Error creating room: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}

	room, err := h.roomUseCase.GetRoom(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) ListByHouseID(c echo.Context) error {
	houseID, err := strconv.ParseInt(c.Param("house_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid house ID format"})
	}

	rooms, err := h.roomUseCase.ListRoomsByHouse(c.Request().Context(), houseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}

	room := new(domain.Room)
	if err := c.Bind(room); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	room.ID = id

	err = h.roomUseCase.UpdateRoom(c.Request().Context(), room)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "room updated successfully"})
}

func (h *RoomHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID format"})
	}

	err = h.roomUseCase.DeleteRoom(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "room deleted successfully"})
}
