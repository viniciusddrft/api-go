package controllers_test

import (
	"api/api/controllers"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
