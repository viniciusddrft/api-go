package main

import (
	routes "api/api/routes"

	db "api/api/db"

	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	app := gin.Default()
	routes.AppRoutes(app)
	app.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Encerrando a conexão com o banco de dados...")
		db.DB.Close()
		os.Exit(0)
	}()
}
