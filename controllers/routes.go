package controllers

import "github.com/nityanandagohain/yield-chaser/utils"

func (s *Server) initializeRoutes() {

	s.Router.HandleFunc("/", utils.SetMiddlewareJSON(s.Home)).Methods("GET")
	s.Router.HandleFunc("/wallet/{id}", utils.SetMiddlewareJSON(s.Read)).Methods("GET")
	s.Router.HandleFunc("/wallet/{id}", utils.SetMiddlewareJSON(s.CreateOrUpdate)).Methods("POST")
	s.Router.HandleFunc("/wallet/{id}", utils.SetMiddlewareJSON(s.Delete)).Methods("DELETE")

}
