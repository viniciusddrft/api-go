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
	err := app.Run()
	if err != nil {
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Encerrando a conexão com o banco de dados...")
		err = db.DB.Close()
		if err != nil {
			return
		}
		os.Exit(0)
	}()
}
