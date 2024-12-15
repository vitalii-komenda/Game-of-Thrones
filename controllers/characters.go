package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vitalii-komenda/got/entities"
)

type CharactersController struct {
	charactersRepo    entities.CharactersRepository
	relationshipsRepo entities.RelationshipsRepository
}

func NewCharactersController(
	charactersRepo entities.CharactersRepository,
	relationshipsRepo entities.RelationshipsRepository,
) *CharactersController {
	return &CharactersController{
		charactersRepo:    charactersRepo,
		relationshipsRepo: relationshipsRepo,
	}
}

// GetAll godoc
// @Summary Get all characters
// @Description Get all characters with pagination
// @Tags characters
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {array} []entities.CharacterEntry
// @Failure 400 {object} map[string]any
// @Router /characters [get]
func (c *CharactersController) GetAll(g *gin.Context) {
	pageStr := g.Query("page")
	if pageStr == "" {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		fmt.Printf("Error parsing page number: %v", err)
		RespondWithError(g, http.StatusBadRequest, err.Error())
		return
	}

	characters, err := c.charactersRepo.GetAll(g.Request.Context(), page)
	if err != nil {
		RespondWithError(g, http.StatusBadRequest, err.Error())
	} else {
		RespondWithJSON(g, http.StatusOK, characters)
	}
}

// Get godoc
// @Summary Get a character by name
// @Description Get a character by name
// @Tags characters
// @Accept  json
// @Produce  json
// @Param name path string true "Character name"
// @Success 200 {object} []entities.CharacterEntry
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /characters/{name} [get]
func (c *CharactersController) Get(g *gin.Context) {
	character := g.Params.ByName("name")
	value, err := c.charactersRepo.Get(g.Request.Context(), character)
	if err != nil {
		RespondWithError(g, http.StatusBadRequest, err.Error())
	} else if len(value) == 0 {
		RespondWithNotFound(g)
	} else {
		RespondWithJSON(g, http.StatusOK, value)
	}
}

// Post godoc
// @Summary Create a new character
// @Description Create a new character and add relationships
// @Tags characters
// @Accept json
// @Produce json
// @Param character body entities.CharacterEntry true "Character Entry"
// @Success 200 {object} entities.CharacterEntry
// @Failure 400 {object} map[string]any
// @Router /characters [post]
func (c *CharactersController) Post(g *gin.Context) {
	var character entities.CharacterEntry
	if err := g.ShouldBindJSON(&character); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.charactersRepo.CreateCharacterAndActor(g.Request.Context(), &character)
	if err == nil {
		err = c.relationshipsRepo.AddAll(g.Request.Context(), character)
	}

	if err != nil {
		RespondWithError(g, http.StatusBadRequest, err.Error())
	} else {
		RespondWithJSON(g, http.StatusOK, character)
	}
}

// Delete godoc
// @Summary Delete a character by name
// @Description Delete a character by name
// @Tags characters
// @Accept  json
// @Produce  json
// @Param name path string true "Character name"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Router /characters/{name} [delete]
func (c *CharactersController) Delete(g *gin.Context) {
	character := g.Params.ByName("name")
	err := c.charactersRepo.Delete(g.Request.Context(), character)

	if err != nil {
		RespondWithError(g, http.StatusBadRequest, err.Error())
	} else {
		RespondWithJSON(g, http.StatusOK, gin.H{})
	}
}

func (c *CharactersController) Put(g *gin.Context) {
	var character entities.CharacterEntry
	if err := g.ShouldBindJSON(&character); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	characterName := g.Params.ByName("name")
	_, err := c.charactersRepo.UpdateCharacterAndActor(g.Request.Context(), &character, characterName)
	if err == nil {
		err = c.relationshipsRepo.UpdateAll(g.Request.Context(), character)
	}
	if err != nil {
		RespondWithError(g, http.StatusBadRequest, err.Error())
	} else {
		RespondWithJSON(g, http.StatusOK, character)
	}
}
