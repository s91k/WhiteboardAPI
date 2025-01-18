package main

import (
	"errors"
	"net/http"
	"regexp"

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

func apiGetBoard(c *gin.Context) {
	enableCors(c)

	var pixels []data.Pixel
	data.DB.Order("y").Order("x").Find(&pixels)

	width := pixels[len(pixels)-1].X + 1
	height := pixels[len(pixels)-1].Y + 1

	grid := make([][]string, height)

	for i := range grid {
		grid[i] = make([]string, width)
	}

	for _, pixel := range pixels {
		grid[pixel.Y][pixel.X] = pixel.Color
	}

	c.JSON(http.StatusOK, grid)
}

func apiGetPixel(c *gin.Context) {
	enableCors(c)

	x := c.Param("x")
	y := c.Param("y")

	var pixel data.Pixel
	err := data.DB.Where("x = ? AND y = ?", x, y).First(&pixel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.IndentedJSON(http.StatusAccepted, pixel)
	}
}

func apiSetPixel(c *gin.Context) {
	enableCors(c)

	x := c.Param("x")
	y := c.Param("y")

	var pixel data.Pixel
	err := data.DB.Where("x = ? AND y = ?", x, y).First(&pixel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		color := c.Query("color")

		if !isColor(color) {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "invalid color"})
		} else {
			data.DB.Where(map[string]interface{}{"x": x, "y": x}).Updates(&data.Pixel{Color: color})

			c.IndentedJSON(http.StatusAccepted, pixel)
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
	router.GET("/board", apiGetBoard)
	router.GET("/pixel/:x/:y", apiGetPixel)
	router.POST("/pixel/:x/:y", apiSetPixel)

	router.Run(":8080")
}
