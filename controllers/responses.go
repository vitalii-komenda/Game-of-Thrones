package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondWithError(g *gin.Context, code int, message string) {
	fmt.Print("error response", message)
	g.JSON(code, gin.H{"error": message})
}

func RespondWithJSON(g *gin.Context, code int, payload interface{}) {
	g.JSON(code, payload)
}

func RespondWithNotFound(g *gin.Context) {
	g.JSON(http.StatusNotFound, gin.H{})
}
