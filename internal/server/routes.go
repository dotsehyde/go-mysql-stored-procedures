package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)
	e.POST("/createTable", s.CreateTableHandler)
	e.GET("/health", s.healthHandler)

	return e
}

func (s *Server) CreateTableHandler(c echo.Context) error {
	_, err := s.db.Db().Exec("CREATE TABLE IF NOT EXISTS messages (id INT AUTO_INCREMENT PRIMARY KEY,content varchar(300) NOT NULL, createdAt TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		return c.JSON(500, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Table created successfully",
	})
}

// Example
func (s *Server) CallStoredProcedure(c echo.Context) error {
	rows, err := s.db.Db().Query("CALL your_stored_procedure(?, ?)", "", "")
	if err != nil {
		return c.JSON(500, map[string]string{
			"message": err.Error(),
		})
	}
	defer rows.Close()

	// Iterate over the result set
	var resultData []map[string]interface{}
	for rows.Next() {
		// Assuming your stored procedure returns columns named "column1" and "column2"
		var column1, column2 string
		if err := rows.Scan(&column1, &column2); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		// Build your result data
		row := map[string]interface{}{
			"column1": column1,
			"column2": column2,
		}
		resultData = append(resultData, row)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "Stored procedure called successfully",
		"resultData": resultData,
	})
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
