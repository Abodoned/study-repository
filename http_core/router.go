package httpcore

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	HTTPHandlers *HTTPHandlers
}

func NewHTTPServer(HTTPHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		HTTPHandlers: HTTPHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()
	router.Path("/miners").Methods("POST").HandlerFunc(s.HTTPHandlers.HandlerBuyMiner)
	router.Path("/miners/cost").Methods("GET").HandlerFunc(s.HTTPHandlers.HandlerGetMinerCost)
	router.Path("/miners").Methods("GET").Queries("class", "{type}").HandlerFunc(s.HTTPHandlers.HandlerGetMinersByClass)
	router.Path("/miners").Methods("GET").HandlerFunc(s.HTTPHandlers.HandlerGetMinersNow)
	router.Path("/equipment/cost").Methods("GET").HandlerFunc(s.HTTPHandlers.HandlerGetEquipmentCost)
	router.Path("/equipment").Methods("POST").HandlerFunc(s.HTTPHandlers.HandlerBuyEquipment)
	router.Path("/equipment").Methods("GET").HandlerFunc(s.HTTPHandlers.HandlerCheckEquipment)
	router.Path("/balance").Methods("GET").HandlerFunc(s.HTTPHandlers.HandlerGetBalance)
	router.Path("/game").Methods("POST").HandlerFunc(s.HTTPHandlers.HandlerStopGame)
	return http.ListenAndServe(":9091", router)
}
