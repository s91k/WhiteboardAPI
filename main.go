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

func isHTMLColor(s string) bool {
	re := regexp.MustCompile("^#[0-9A-Fa-f]{6}$")
	return re.MatchString(s)
}

func apiGetBoard(c *gin.Context) {
	enableCors(c)

	var pixels []data.Pixel
	data.DB.Order("Y").Order("X").Find(&pixels)

	c.JSON(http.StatusOK, pixels)
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

		if !isHTMLColor(color) {
			c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "invalid color"})
		} else {
			pixel.Color = color
			data.DB.Save(&pixel)

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
	router.GET("/api/board", apiGetBoard)
	router.GET("/api/pixel/:x/:y", apiGetPixel)
	router.POST("/api/pixel/:x/:y", apiSetPixel)

	router.Run(":8080")
}
