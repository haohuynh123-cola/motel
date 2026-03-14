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
	contracts.GET("", h.ListAll)
	contracts.GET("/:id", h.GetByID)

	g.GET("/houses/:house_id/contracts", h.ListByHouseID)
}

func (h *ContractHandler) Create(c echo.Context) error {
	contract := new(domain.Contract)
	if err := c.Bind(contract); err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "Dữ liệu hợp đồng không hợp lệ"})
	}

	if err := h.contractUC.CreateContract(c.Request().Context(), contract); err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: err.Error()})
	}

	return c.JSON(http.StatusCreated, port.ApiResponse{Status: true, Data: contract})
}

func (h *ContractHandler) ListAll(c echo.Context) error {
	response, err := h.contractUC.ListAllContracts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *ContractHandler) GetByID(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	response, err := h.contractUC.GetContract(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, port.ApiResponse{Status: false, Data: "Không tìm thấy hợp đồng"})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *ContractHandler) ListByHouseID(c echo.Context) error {
	houseID, err := strconv.ParseInt(c.Param("house_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, port.ApiResponse{Status: false, Data: "invalid house ID format"})
	}

	response, err := h.contractUC.ListContractsByHouse(c.Request().Context(), houseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, port.ApiResponse{Status: false, Data: "internal server error"})
	}

	return c.JSON(http.StatusOK, response)
}
