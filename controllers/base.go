package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/nityanandagohain/yield-chaser/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbUrl, addr string) {

	var err error

	server.DB, err = gorm.Open("postgres", DbUrl)
	if err != nil {
		fmt.Println(err)
	}

	err = server.DB.DB().Ping()
	if err != nil {
		log.Fatalf(err.Error())
	}

	server.DB.Debug().AutoMigrate(&models.YChaseModel{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()

	handler := cors.AllowAll().Handler(server.Router)

	fmt.Println("Listening to port", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
