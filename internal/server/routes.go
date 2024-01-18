package server

import (
	"net/http"
	"store-procedures-mysql/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)
	e.POST("/createTable", s.CreateTableHandler)
	e.POST("/create", s.CreateMessageHandler)
	e.PUT("/update/:id", s.UpdateMessageHandler)
	e.GET("/get/:id", s.GetMessageByIDHandler)
	e.DELETE("/delete/:id", s.DeleteMessageHandler)
	e.GET("/all", s.GetAllMessageHandler)
	e.GET("/health", s.healthHandler)

	return e
}

func (s *Server) CreateTableHandler(c echo.Context) error {
	_, err := s.db.Db().Exec("CREATE TABLE IF NOT EXISTS messages (id INT AUTO_INCREMENT PRIMARY KEY, content VARCHAR(300) NOT NULL,createdAt TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6));")
	if err != nil {
		c.JSON(500, map[string]string{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Table created successfully",
	})
}

// Create Message
func (s *Server) CreateMessageHandler(c echo.Context) error {
	var msg model.Message
	if err := c.Bind(&msg); err != nil {
		c.JSON(400, map[string]any{
			"message": err.Error(),
		})
		return echo.ErrBadGateway
	}
	res, err := s.db.Db().Query("CALL messages_create(?)", msg.Content)
	if err != nil {
		c.JSON(500, map[string]any{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	for res.Next() {
		if err := res.Scan(&msg.ID, &msg.Content, &msg.CreatedAt); err != nil {
			c.JSON(500, map[string]any{
				"message": err.Error(),
			})
			return echo.ErrInternalServerError
		}
	}
	return c.JSON(200, map[string]any{
		"message": "Message created",
		"data":    msg,
	})
}

// Get All Messages
func (s *Server) GetAllMessageHandler(c echo.Context) error {
	rows, err := s.db.Db().Query("CALL messages_read_all")
	var resultData []model.Message
	if err != nil {
		c.JSON(500, map[string]string{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		// Assuming your stored procedure returns columns named "column1" and "column2"
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		// Build your result data
		resultData = append(resultData, msg)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": resultData,
	})
}

// Get Message by ID
func (s *Server) GetMessageByIDHandler(c echo.Context) error {
	id := c.Param("id")
	rows, err := s.db.Db().Query("CALL messages_read_by_id(?)", id)
	if err != nil {
		c.JSON(500, map[string]any{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	var msg model.Message
	for rows.Next() {
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt); err != nil {
			c.JSON(500, map[string]any{
				"message": err.Error(),
			})
			return echo.ErrInternalServerError
		}
	}
	if msg.ID == 0 {
		return c.JSON(200, map[string]any{
			"data": map[string]any{},
		})
	}
	return c.JSON(200, map[string]any{
		"data": msg,
	})
}

// Update Message
func (s *Server) UpdateMessageHandler(c echo.Context) error {
	var msg model.Message
	if err := c.Bind(&msg); err != nil {
		c.JSON(500, map[string]any{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	id := c.Param("id")
	_, err := s.db.Db().Exec("CALL messages_update(?,?)", id, msg.Content)
	if err != nil {
		c.JSON(500, map[string]any{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	return c.JSON(200, map[string]any{
		"message": "Message update",
	})

}

// Delete Message
func (s *Server) DeleteMessageHandler(c echo.Context) error {
	id := c.Param("id")
	_, err := s.db.Db().Exec("CALL messages_delete(?)", id)
	if err != nil {
		c.JSON(500, map[string]any{
			"message": err.Error(),
		})
		return echo.ErrInternalServerError
	}
	return c.JSON(200, map[string]any{
		"message": "Message deleted",
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
