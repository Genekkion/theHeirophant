package main

import (
	"encoding/json"
	"io"
	"net/http"

	"theHeirophant/log"

	"github.com/gorilla/mux"
)

type Server struct {
	router mux.Router
	port   string
}

func initServer(port string) Server {

	server := Server{
		router: *mux.NewRouter(),
		port:   port,
	}
	server.addRoutes()
	return server
}
func (server *Server) addRoutes() {
	server.router.HandleFunc("/health", server.checkHealth)
	server.router.HandleFunc("/echo", server.echoHandler)
	server.router.HandleFunc("/do", server.doRequest)
	server.router.NotFoundHandler = NotFoundHandler{}
	log.Info("Routes added")
}

func (server Server) run() {
	address := "http://localhost:" + server.port
	log.Info("Server starting at: " + address)
	http.ListenAndServe("localhost:"+server.port, &server.router)
}

func (server Server) checkHealth(writer http.ResponseWriter,
	request *http.Request) {
	log.Info("/ called", "header", request.Header, "body", request.Body)
	writeJson(writer, http.StatusOK, map[string]interface{}{
		"message": "The Heirophant is up and running.",
	})
}

func (server Server) echoHandler(writer http.ResponseWriter,
	request *http.Request) {
	log.Info("/echo called", "header", request.Header, "body", request.Body)

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		errorResponse(writer, request, err)
		return
	}

	writeJson(writer, http.StatusOK, map[string]interface{}{
		"message": string(requestBody),
	})
}

type NotFoundHandler struct {
}

func (handler NotFoundHandler) ServeHTTP(writer http.ResponseWriter,
	request *http.Request) {
	log.Error("Unregistered endpoint called.", request.URL)
	writeJson(writer, http.StatusNotFound, map[string]interface{}{
		"error": "Endpoint not found.",
	})
}

func writeJson(writer http.ResponseWriter, statusCode int, value any) error {
	writer.WriteHeader(statusCode)
	writer.Header().Add("Content-Type", "application/json")

	return log.ErrorWrapper(json.NewEncoder(writer).Encode(value))
}
