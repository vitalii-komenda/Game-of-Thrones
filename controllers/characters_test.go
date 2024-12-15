package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vitalii-komenda/got/entities"
	"github.com/vitalii-komenda/got/mocks"
)

func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCharactersRepo := new(mocks.CharactersRepositoryMock)
	mockRelationshipsRepo := new(mocks.RelationshipsRepositoryMock)
	controller := NewCharactersController(mockCharactersRepo, mockRelationshipsRepo)

	t.Run("success", func(t *testing.T) {
		mockCharactersRepo.GetAllFunc = func(ctx context.Context, page int) ([]entities.CharacterEntry, error) {
			return []entities.CharacterEntry{{CharacterName: "test"}}, nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/characters?page=0", nil)

		controller.GetAll(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed to getAll", func(t *testing.T) {
		mockCharactersRepo.GetAllFunc = func(ctx context.Context, page int) ([]entities.CharacterEntry, error) {
			return nil, fmt.Errorf("some error")
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/characters?page=0", nil)

		controller.GetAll(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("no page passed", func(t *testing.T) {
		mockCharactersRepo.GetAllFunc = func(ctx context.Context, page int) ([]entities.CharacterEntry, error) {
			return []entities.CharacterEntry{{CharacterName: "test"}}, nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/characters", nil)

		controller.GetAll(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid page", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/characters?page=invalid", nil)

		controller.GetAll(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func TestGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCharactersRepo := new(mocks.CharactersRepositoryMock)
	mockRelationshipsRepo := new(mocks.RelationshipsRepositoryMock)
	controller := NewCharactersController(mockCharactersRepo, mockRelationshipsRepo)

	t.Run("success", func(t *testing.T) {
		mockCharactersRepo.GetFunc = func(ctx context.Context, name string) ([]entities.CharacterEntry, error) {
			return []entities.CharacterEntry{{CharacterName: "test"}}, nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "name", Value: "test"}}
		c.Request, _ = http.NewRequest("GET", "/characters/test", nil)

		controller.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("failed to get", func(t *testing.T) {
		mockCharactersRepo.GetFunc = func(ctx context.Context, name string) ([]entities.CharacterEntry, error) {
			return nil, fmt.Errorf("some error")
		}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "name", Value: "test"}}
		c.Request, _ = http.NewRequest("GET", "/characters/test", nil)

		controller.Get(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockCharactersRepo.GetFunc = func(ctx context.Context, name string) ([]entities.CharacterEntry, error) {
			return nil, nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "name", Value: "test"}}
		c.Request, _ = http.NewRequest("GET", "/characters/test", nil)

		controller.Get(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCharactersRepo := new(mocks.CharactersRepositoryMock)
	mockRelationshipsRepo := new(mocks.RelationshipsRepositoryMock)
	controller := NewCharactersController(mockCharactersRepo, mockRelationshipsRepo)

	t.Run("success", func(t *testing.T) {
		mockCharactersRepo.CreateCharacterAndActorFunc = func(ctx context.Context, character *entities.CharacterEntry) error {
			return nil
		}
		mockRelationshipsRepo.AddAllFunc = func(ctx context.Context, character entities.CharacterEntry) error {
			return nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/characters", strings.NewReader(`{"name":"test"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Post(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/characters", strings.NewReader(`invalid`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Post(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCharactersRepo := new(mocks.CharactersRepositoryMock)
	mockRelationshipsRepo := new(mocks.RelationshipsRepositoryMock)
	controller := NewCharactersController(mockCharactersRepo, mockRelationshipsRepo)

	t.Run("success", func(t *testing.T) {
		mockCharactersRepo.DeleteFunc = func(ctx context.Context, name string) error {
			return nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "name", Value: "test"}}
		c.Request, _ = http.NewRequest("DELETE", "/characters/test", nil)

		controller.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		mockCharactersRepo.DeleteFunc = func(ctx context.Context, name string) error {
			return fmt.Errorf("some error")
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "name", Value: "test"}}
		c.Request, _ = http.NewRequest("DELETE", "/characters/test", nil)

		controller.Delete(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPut(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCharactersRepo := new(mocks.CharactersRepositoryMock)
	mockRelationshipsRepo := new(mocks.RelationshipsRepositoryMock)
	controller := NewCharactersController(mockCharactersRepo, mockRelationshipsRepo)

	t.Run("success", func(t *testing.T) {
		mockCharactersRepo.UpdateCharacterAndActorFunc = func(ctx context.Context, characterEntryEntry *entities.CharacterEntry, characterName string) (int, error) {
			return 1, nil
		}
		mockRelationshipsRepo.UpdateAllFunc = func(ctx context.Context, character entities.CharacterEntry) error {
			return nil
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "name", Value: "test"}}
		c.Request, _ = http.NewRequest("PUT", "/characters/test", strings.NewReader(`{"name":"test"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Put(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/characters/test", strings.NewReader(`invalid`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Put(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}
