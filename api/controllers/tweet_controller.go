package controllers

import (
	entities "api/api/entities"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TweetController struct {
	db *sql.DB
}

func NewTweetController(db *sql.DB) *TweetController {
	return &TweetController{db: db}
}

func (t *TweetController) FindAll(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodGet {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Método não aceito"})
		return
	}

	rows, err := t.db.Query("SELECT * FROM tweets")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar tweets: " + err.Error()})
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var tweets []entities.Tweet

	for rows.Next() {
		var tweet entities.Tweet
		err := rows.Scan(&tweet.Id, &tweet.Description)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar resultado: " + err.Error()})
			return
		}
		tweets = append(tweets, tweet)
	}

	if err = rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro na iteração: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tweets)
}

func (t *TweetController) Create(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodPost {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Método não aceito"})
		return
	}

	var tweet entities.Tweet

	if err := ctx.BindJSON(&tweet); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var id int64
	err := t.db.QueryRow("INSERT INTO tweets (description) VALUES ($1) RETURNING id", tweet.Description).Scan(&id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar tweet: " + err.Error()})
		return
	}

	tweet.Id = id

	ctx.JSON(http.StatusOK, tweet)
}

func (t *TweetController) Update(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodPut {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Método não aceito"})
		return
	}

	id := ctx.Param("id")

	var tweet entities.Tweet

	if err := ctx.BindJSON(&tweet); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	result, err := t.db.Exec("UPDATE tweets SET description = $1 WHERE id = $2", tweet.Description, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar tweet: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar linhas afetadas: " + err.Error()})
		return
	}

	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tweet não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tweet atualizado com sucesso"})
}

func (t *TweetController) Delete(ctx *gin.Context) {
	if ctx.Request.Method != http.MethodDelete {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Método não aceito"})
		return
	}

	id := ctx.Param("id")

	result, err := t.db.Exec("DELETE FROM tweets WHERE id = $1", id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir tweet: " + err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	fmt.Println(rowsAffected)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar linhas afetadas: " + err.Error()})
		return
	}
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tweet não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tweet excluído com sucesso"})
}
