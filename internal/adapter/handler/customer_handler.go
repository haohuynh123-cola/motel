package handler

import (
	"net/http"
	"strconv"
	"tro-go/internal/domain"
	"tro-go/internal/port"

	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	customerUC port.CustomerUseCase
}

func NewCustomerHandler(g *echo.Group, customerUC port.CustomerUseCase) {
	h := &CustomerHandler{customerUC: customerUC}

	customers := g.Group("/customers")
	customers.POST("", h.Register)
	customers.GET("", h.List)
	customers.GET("/:id", h.GetByID)
	customers.PUT("/:id", h.Update)
	customers.DELETE("/:id", h.Delete)
}

func (h *CustomerHandler) Register(c echo.Context) error {
	customer := new(domain.Customer)
	if err := c.Bind(customer); err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "Dữ liệu không hợp lệ"})
	}

	if err := h.customerUC.RegisterCustomer(c.Request().Context(), customer); err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}

	return c.JSON(http.StatusCreated, port.ApiResponse{Status: true, Data: customer})
}

func (h *CustomerHandler) GetByID(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	response, err := h.customerUC.GetCustomer(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, port.ApiResponse{Status: false, Data: "Không tìm thấy khách thuê"})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CustomerHandler) List(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if parsed, err := strconv.Atoi(pageStr); err == nil {
			page = parsed
		}
	}
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	response, err := h.customerUC.ListCustomers(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CustomerHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	customer := new(domain.Customer)
	if err := c.Bind(customer); err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "Dữ liệu không hợp lệ"})
	}
	customer.ID = id
	err := h.customerUC.UpdateCustomer(c.Request().Context(), customer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}
	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: "Cập nhật thành công"})
}

func (h *CustomerHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := h.customerUC.DeleteCustomer(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}
	return c.JSON(http.StatusOK, port.ApiResponse{Status: true, Data: "Xóa thành công"})
}
