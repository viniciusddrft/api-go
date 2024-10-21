package controllers_test

import (
	"api/api/controllers"
	"api/api/entities"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	controller := controllers.NewTweetController(db)
	rows := sqlmock.NewRows([]string{"id", "description"}).
		AddRow(1, "Primeiro tweet").
		AddRow(2, "Segundo tweet")
	mock.ExpectQuery("SELECT \\* FROM tweets").WillReturnRows(rows)
	req, err := http.NewRequest(http.MethodGet, "/tweets", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	controller.FindAll(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `[{"id":1,"description":"Primeiro tweet"},{"id":2,"description":"Segundo tweet"}]`
	assert.JSONEq(t, expectedResponse, w.Body.String())
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	controller := controllers.NewTweetController(db)

	newTweet := entities.Tweet{
		Description: "Novo tweet de teste",
	}

	mock.ExpectQuery("INSERT INTO tweets").WithArgs(newTweet.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	tweetJSON, _ := json.Marshal(newTweet)
	req, _ := http.NewRequest(http.MethodPost, "/tweets", bytes.NewBuffer(tweetJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	controller.Create(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var responseTweet entities.Tweet
	err = json.Unmarshal(w.Body.Bytes(), &responseTweet)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), responseTweet.Id)
	assert.Equal(t, newTweet.Description, responseTweet.Description)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	controller := controllers.NewTweetController(db)

	updatedTweet := entities.Tweet{
		Description: "Tweet atualizado",
	}

	mock.ExpectExec("UPDATE tweets SET description").WithArgs(updatedTweet.Description, "1").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simular que uma linha foi afetada

	tweetJSON, _ := json.Marshal(updatedTweet)
	req, _ := http.NewRequest(http.MethodPut, "/tweets/1", bytes.NewBuffer(tweetJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	controller.Update(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var jsonResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Tweet atualizado com sucesso", jsonResponse["message"])
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	controller := controllers.NewTweetController(db)

	mock.ExpectExec("DELETE FROM tweets WHERE id = \\$1"). // Usar \\$1 para escapar no regex
								WithArgs("1").                            // Simulando que estamos tentando excluir o tweet com ID 1
								WillReturnResult(sqlmock.NewResult(0, 1)) // Simular que uma linha foi afetada

	req, _ := http.NewRequest(http.MethodDelete, "/tweets/1", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}} // Simulando que a rota tem um parâmetro "id" e
	controller.Delete(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var jsonResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)
	assert.Equal(t, "Tweet excluído com sucesso", jsonResponse["message"])
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
