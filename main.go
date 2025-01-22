package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/s91k/WhiteboardAPI/data"
	"gorm.io/gorm"
)

func enableCors(c *gin.Context) {
	(*c).Header("Access-Control-Allow-Origin", "*")
}

func start(c *gin.Context) {
	c.Status(http.StatusOK)
}

func isColor(s string) bool {
	re := regexp.MustCompile("^[0-9A-Fa-f]{6}$")
	return re.MatchString(s)
}

func apiGetBoards(c *gin.Context) {
	enableCors(c)

	var whiteboards []data.Whiteboard
	data.DB.Find(&whiteboards)

	c.JSON(http.StatusOK, whiteboards)
}

func apiGetBoard(c *gin.Context) {
	enableCors(c)

	i := c.Param("i")

	var whiteboard data.Whiteboard
	data.DB.Where(map[string]interface{}{"id": i}).First(&whiteboard)

	var pixels []data.Pixel
	data.DB.Where(map[string]interface{}{"whiteboard_id": i}).Find(&pixels)

	grid := make([][]string, whiteboard.Height)

	for i := range grid {
		grid[i] = make([]string, whiteboard.Width)

		for j := range grid[i] {
			grid[i][j] = whiteboard.DefaultColor
		}
	}

	for _, pixel := range pixels {
		grid[pixel.Y][pixel.X] = pixel.Color
	}

	c.JSON(http.StatusOK, grid)
}

func apiCreateBoard(c *gin.Context) {
	enableCors(c)

	var whiteboard data.Whiteboard
	err := c.BindJSON(&whiteboard)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
	} else {
		data.DB.Create(&whiteboard)

		c.IndentedJSON(http.StatusCreated, whiteboard)
	}
}

func apiUpdateBoard(c *gin.Context) {
	enableCors(c)

	i := c.Param("i")

	var whiteboard data.Whiteboard
	err := data.DB.Where(map[string]interface{}{"id": i}).First(&whiteboard).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		err := c.BindJSON(&whiteboard)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		} else {
			data.DB.Where(map[string]interface{}{"id": i}).Save(&whiteboard)

			data.DB.Where(map[string]interface{}{"whiteboard_id": i}).Where("x >= ?", whiteboard.Width).Delete(&data.Pixel{})
			data.DB.Where(map[string]interface{}{"whiteboard_id": i}).Where("y >= ?", whiteboard.Height).Delete(&data.Pixel{})

			c.IndentedJSON(http.StatusAccepted, whiteboard)
		}
	}
}

func apiDeleteBoard(c *gin.Context) {
	enableCors(c)

	i := c.Param("i")

	var whiteboard data.Whiteboard
	err := data.DB.Where(map[string]interface{}{"id": i}).First(&whiteboard).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		data.DB.Where(map[string]interface{}{"whiteboard_id": i}).Delete(&data.Pixel{})
		data.DB.Where(map[string]interface{}{"id": i}).Delete(&whiteboard)

		c.IndentedJSON(http.StatusAccepted, whiteboard)
	}
}

func apiGetPixel(c *gin.Context) {
	enableCors(c)

	i := c.Param("i")
	x := c.Param("x")
	y := c.Param("y")

	var pixel data.Pixel
	err := data.DB.Where(map[string]interface{}{"whiteboard_id": i, "x": x, "y": y}).First(&pixel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.IndentedJSON(http.StatusAccepted, pixel)
	}
}

func apiSetPixel(c *gin.Context) {
	enableCors(c)

	i := c.Param("i")
	xStr := c.Param("x")
	yStr := c.Param("y")

	x, err := strconv.Atoi(xStr)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
	}

	y, err := strconv.Atoi(yStr)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
	}

	var whiteboard data.Whiteboard
	err = data.DB.Where(map[string]interface{}{"id": i}).First(&whiteboard).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		if x < 0 || x >= whiteboard.Width || y < 0 || y >= whiteboard.Height {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "out of bounds"})
		} else {
			color := c.Query("color")

			if !isColor(color) {
				c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "invalid color"})
			} else {
				var pixel data.Pixel
				data.DB.Where(map[string]interface{}{"whiteboard_id": i, "x": x, "y": y}).FirstOrCreate(&pixel)
				pixel.Color = color

				data.DB.Where(map[string]interface{}{"whiteboard_id": i, "x": x, "y": y}).Updates(&pixel)

				c.IndentedJSON(http.StatusAccepted, pixel)
			}
		}
	}
}

var config Config

func main() {
	readConfig(&config)

	data.InitDatabase(config.Database.File,
		config.Database.Server,
		config.Database.Database,
		config.Database.Username,
		config.Database.Password,
		config.Database.Port,
		config.Whiteboard.Width,
		config.Whiteboard.Height)

	router := gin.Default()

	router.GET("/", start)
	router.GET("/api/board", apiGetBoards)
	router.GET("/api/board/:i", apiGetBoard)
	router.POST("/api/board", apiCreateBoard)
	router.PUT("/api/board/:i", apiUpdateBoard)
	router.DELETE("/api/board/:i", apiDeleteBoard)
	router.GET("/api/board/:i/pixel/:x/:y", apiGetPixel)
	router.POST("/api/board/:i/pixel/:x/:y", apiSetPixel)

	router.Run(":8080")
}
