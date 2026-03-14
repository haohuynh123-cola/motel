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
	e.POST("/rooms/:id/remind", handler.Remind)
	e.POST("/rooms/:id/book", handler.BookAppointment)
}

func (h *RoomHandler) Remind(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid ID format"})
	}

	var req struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&req); err != nil || req.Email == "" {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "vui lòng cung cấp email người nhận"})
	}

	err = h.roomUseCase.SendPaymentReminder(c.Request().Context(), id, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}

	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: "đã gửi email nhắc nhở thành công"})
}

func (h *RoomHandler) BookAppointment(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid ID format"})
	}

	app := new(domain.Appointment)
	if err := c.Bind(app); err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: err.Error()})
	}
	app.RoomID = id

	err = h.roomUseCase.BookAppointment(c.Request().Context(), app)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}

	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: "đã đặt lịch hẹn xem phòng thành công, vui lòng kiểm tra email của bạn"})
}

func (h *RoomHandler) Create(c echo.Context) error {
	room := new(domain.Room)
	if err := c.Bind(room); err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: err.Error()})
	}

	err := h.roomUseCase.CreateRoom(c.Request().Context(), room)
	if err != nil {
		log.Printf("Error creating room: %v\n", err)
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}

	return c.JSON(http.StatusCreated, port.ApiResponse{Status: true, Data: room})
}

func (h *RoomHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid ID format"})
	}

	room, err := h.roomUseCase.GetRoom(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, port.ApiResponse{Status: false, Data: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}

	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: room})
}

func (h *RoomHandler) ListByHouseID(c echo.Context) error {
	houseID, err := strconv.ParseInt(c.Param("house_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid house ID format"})
	}

	response, err := h.roomUseCase.ListRoomsByHouse(c.Request().Context(), houseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *RoomHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid ID format"})
	}

	room := new(domain.Room)
	if err := c.Bind(room); err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: err.Error()})
	}
	room.ID = id

	err = h.roomUseCase.UpdateRoom(c.Request().Context(), room)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, port.ApiResponse{Status: false, Data: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}

	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: "room updated successfully"})
}

func (h *RoomHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid ID format"})
	}

	err = h.roomUseCase.DeleteRoom(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, port.ApiResponse{Status: false, Data: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}

	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: "room deleted successfully"})
}
