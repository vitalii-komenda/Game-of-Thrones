package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vitalii-komenda/got/entities"
	"github.com/vitalii-komenda/got/services"
)

type SearchController struct {
	charactersRepo entities.CharactersRepository
}

func NewSearchController(
	characterRepo entities.CharactersRepository,
) *SearchController {
	return &SearchController{
		charactersRepo: characterRepo,
	}
}

// GetFromElastic godoc
// @Summary Search characters in elastic
// @Description Search characters by term in elastic
// @Tags search
// @Accept  json
// @Produce  json
// @Param term query string true "Search term"
// @Success 200 {array} []entities.CharacterEntry
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /elastic/search [get]
func (c *SearchController) GetFromElastic(g *gin.Context) {
	term := g.Query("term")
	value, err := services.SendElasticSearchRequest(term)

	if err != nil {
		RespondWithError(g, http.StatusBadRequest, err.Error())
	} else if len(value) == 0 {
		RespondWithNotFound(g)
	} else {
		RespondWithJSON(g, http.StatusOK, value)
	}
}
