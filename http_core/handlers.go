package httpcore

import (
	"MinerGame/repositor"
	"encoding/json"
	"net/http"
)

type HTTPHandlers struct {
	RepositoryFactory *repositor.Factory
}

func NewHTTPHandlers(RepositoryFactory *repositor.Factory) *HTTPHandlers {
	return &HTTPHandlers{
		RepositoryFactory: RepositoryFactory,
	}
}

func (h *HTTPHandlers) HandlerBuyMiner(w http.ResponseWriter, r *http.Request) {
	/*
	   API Rest
	   Description:
	   Wait JSON
	   Pattern /miners
	   method Post
	   Status code 201 succesed
	   return JSON with miner info? не увеерен что правильно?

	   	Failed  400 Bad Request-----409 Conflict------500 Internal Server Error
	   	return JSON Err+Time
	*/
	var minerClass BuyAnyDTO
	if err := json.NewDecoder(r.Body).Decode(&minerClass); err != nil {
		WriteError(w, err, 400)
		return
	}
	minerInfo, err := h.RepositoryFactory.BuyMiner(minerClass.Name)
	if err != nil {
		SelectError(w, err)
		return
	}

	WriteResponse(w, minerInfo, 201)

}

func (h *HTTPHandlers) HandlerGetMinerCost(w http.ResponseWriter, r *http.Request) {
	/*
	   Pattern /miners/cost
	   Method GET
	   response: Successed: JSON+Status code 200
	   failed: Console response
	*/
	WriteResponse(w, h.RepositoryFactory.GetMinerCost(), http.StatusOK)
}

func (h *HTTPHandlers) HandlerGetMinersNow(w http.ResponseWriter, r *http.Request) {
	/*
	   Pattern /miners
	   Method Get
	   response: Successed: JSON+Status code 200
	   failed: console response
	*/
	WriteResponse(w, h.RepositoryFactory.GetMinersNow(), http.StatusOK)
}
func (h *HTTPHandlers) HandlerGetMinersByClass(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	class := query.Get("class")
	WriteResponse(w, h.RepositoryFactory.GetMinerByClass(class), http.StatusOK)
}
func (h *HTTPHandlers) HandlerGetEquipmentCost(w http.ResponseWriter, r *http.Request) {
	/*
	   Pattern /Equipment/cost
	   Method GET
	   response: Successed: JSON+Status code 200
	   failed: JSON ERROR+Time  Status Code ----500 Internal Server Error
	*/
	WriteResponse(w, h.RepositoryFactory.GetEquipmentCost(), 200)
}
func (h *HTTPHandlers) HandlerBuyEquipment(w http.ResponseWriter, r *http.Request) {
	/*
	   API Rest
	   Description:
	   Wait JSON
	   Pattern /equipment
	   method Post
	   Status code 201 succesed
	   return JSON with Equipment

	   	Failed  400 Bad Request-----409 Conflict------500 Internal Server Error
	   	return JSON Err+Time
	*/
	var equipmentDTO BuyAnyDTO

	if err := json.NewDecoder(r.Body).Decode(&equipmentDTO); err != nil {
		WriteError(w, err, 400)
		return
	}
	equipment, err := h.RepositoryFactory.BuyEquipment(equipmentDTO.Name)
	if err != nil {
		SelectError(w, err)
		return
	}
	WriteResponse(w, equipment, http.StatusCreated)

}

func (h *HTTPHandlers) HandlerGetBalance(w http.ResponseWriter, r *http.Request) {
	/*
		   Pattern /balance
		   Method Get
		response: Successed: JSON+Status code 200
		failed: console response
	*/
	WriteResponse(w, h.RepositoryFactory.CheckBalance(), http.StatusOK)
}

func (h *HTTPHandlers) HandlerCheckEquipment(w http.ResponseWriter, r *http.Request) {
	/*
	   Pattern /equipment
	   Method Get
	   Response JSON+Status code 200
	   Failed: console response
	*/
	WriteResponse(w, h.RepositoryFactory.CheckEquipment(), http.StatusOK)
}

func (h *HTTPHandlers) HandlerStopGame(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, h.RepositoryFactory.StopGame(), http.StatusOK)
}
