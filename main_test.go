package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	router := gin.Default()
	router.GET("/", start)

	req, err := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestCorrectColorIsAllowed(t *testing.T) {
	color := "FF9933"

	isColor := isColor(color)

	assert.True(t, isColor)
}

func TestColorWithHashIsNotAllowed(t *testing.T) {
	color := "#FF9933"

	isColor := isColor(color)

	assert.False(t, isColor)
}

func TestColorWithIncorrectValuesNotAllowed(t *testing.T) {
	color := "UU99ÅÅ"

	isColor := isColor(color)

	assert.False(t, isColor)
}
