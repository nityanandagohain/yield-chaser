package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nityanandagohain/yield-chaser/models"
	"github.com/nityanandagohain/yield-chaser/utils"
)

// func (server *Server) Create(w http.ResponseWriter, r *http.Request) {

// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		utils.ERROR(w, http.StatusUnprocessableEntity, err)
// 	}
// 	user := models.YChaseModel{}
// 	err = json.Unmarshal(body, &user)
// 	if err != nil {
// 		utils.ERROR(w, http.StatusUnprocessableEntity, err)
// 		return
// 	}

// 	entryCreated, err := user.Create(server.DB)

// 	if err != nil {

// 		formattedError := utils.FormatError(err.Error())

// 		utils.ERROR(w, http.StatusInternalServerError, formattedError)
// 		return
// 	}
// 	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, entryCreated.WalletAddress))
// 	utils.JSON(w, http.StatusCreated, entryCreated)
// }

func (server *Server) Read(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	walletAddress := vars["id"]

	user := models.YChaseModel{}
	userGotten, err := user.Read(server.DB, walletAddress)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	utils.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	walletAddress := vars["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.YChaseModel{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	readUser := models.YChaseModel{}
	// check if present
	_, err = readUser.Read(server.DB, walletAddress)
	if err != nil {
		newUser := models.YChaseModel{
			WalletAddress: walletAddress,
			Assets:        user.Assets,
		}
		entryCreated, err := newUser.Create(server.DB)
		if err != nil {
			formattedError := utils.FormatError(err.Error())
			utils.ERROR(w, http.StatusInternalServerError, formattedError)
			return
		}
		utils.JSON(w, http.StatusOK, entryCreated)
		return
	}

	updatedModel, err := user.Update(server.DB, string(walletAddress))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	updatedModel.WalletAddress = walletAddress
	utils.JSON(w, http.StatusOK, updatedModel)
}

func (server *Server) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.YChaseModel{}

	walletAddress := vars["id"]

	_, err := user.Delete(server.DB, walletAddress)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", walletAddress))
	utils.JSON(w, http.StatusNoContent, "")
}
