package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewDocsHandler(e *echo.Echo) {
	// Serve the static openapi.yaml file
	e.File("/docs/openapi.yaml", "openapi.yaml")

	// Render the ReDoc HTML page
	e.GET("/docs", func(c echo.Context) error {
		html := `
<!DOCTYPE html>
<html>
  <head>
    <title>Tro-Go API Documentation</title>
    <!-- needed for adaptive design -->
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

    <!--
    ReDoc doesn't change outer page styles
    -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='/docs/openapi.yaml'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
  </body>
</html>
`
		return c.HTML(http.StatusOK, html)
	})
}
