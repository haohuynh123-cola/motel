package handler

import (
	"net/http"
	"strconv"
	"tro-go/internal/domain"
	"tro-go/internal/port"

	"github.com/labstack/echo/v4"
)

type ContractHandler struct {
	contractUC port.ContractUseCase
}

func NewContractHandler(g *echo.Group, contractUC port.ContractUseCase) {
	h := &ContractHandler{contractUC: contractUC}

	contracts := g.Group("/contracts")
	contracts.POST("", h.Create)
	contracts.GET("/:id", h.GetByID)
}

func (h *ContractHandler) Create(c echo.Context) error {
	contract := new(domain.Contract)
	if err := c.Bind(contract); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Dữ liệu hợp đồng không hợp lệ"})
	}

	if err := h.contractUC.CreateContract(c.Request().Context(), contract); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, contract)
}

func (h *ContractHandler) GetByID(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	contract, err := h.contractUC.GetContract(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Không tìm thấy hợp đồng"})
	}
	return c.JSON(http.StatusOK, contract)
}
