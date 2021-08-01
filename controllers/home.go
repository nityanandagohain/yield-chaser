package controllers

import (
	"net/http"

	"github.com/nityanandagohain/yield-chaser/utils"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, "Yield Chaser API")

}
