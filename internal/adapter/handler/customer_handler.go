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
}

func (h *CustomerHandler) Register(c echo.Context) error {
	customer := new(domain.Customer)
	if err := c.Bind(customer); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dữ liệu không hợp lệ"})
	}

	if err := h.customerUC.RegisterCustomer(c.Request().Context(), customer); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) GetByID(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	customer, err := h.customerUC.GetCustomer(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Không tìm thấy khách thuê"})
	}
	return c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) List(c echo.Context) error {
	customers, err := h.customerUC.ListCustomers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, customers)
}
