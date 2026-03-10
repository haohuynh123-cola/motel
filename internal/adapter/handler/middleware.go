package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// PermissionMiddleware kiểm tra xem user có quyền cụ thể không
func PermissionMiddleware(requiredPermission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Lấy token từ context
			userToken, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Không tìm thấy thông tin đăng nhập!"})
			}

			// 2. Trích xuất claims
			claims := userToken.Claims.(*jwt.MapClaims)
			
			// 3. Lấy mảng permissions
			permissionsInterface, ok := (*claims)["permissions"].([]interface{})
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Không tìm thấy danh sách quyền!"})
			}

			// 4. Kiểm tra quyền
			hasPerm := false
			for _, p := range permissionsInterface {
				if str, ok := p.(string); ok && str == requiredPermission {
					hasPerm = true
					break
				}
			}

			if !hasPerm {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Bạn không có quyền thực hiện hành động này!"})
			}

			return next(c)
		}
	}
}
