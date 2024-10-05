package user

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type API struct {
	service *Service
}

func NewApi(service *Service) *API {
	return &API{service: service}
}

func sendInternalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	createDTO := &CreateDTO{}
	if err := json.NewDecoder(r.Body).Decode(createDTO); err != nil {
		sendInternalServerError(w)
		return
	}
	newUser, err := a.service.Save(createDTO)
	if err != nil {
		sendInternalServerError(w)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newUser)
	//err = json.NewEncoder(w).Encode(newUser)
	//if err != nil {
	//	sendInternalServerError(w)
	//	return
	//}
}

func (a *API) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := a.service.FindOne(id)
	if err != nil {
		//render.NoContent(w, r)
		render.Status(r, http.StatusNotFound)
		//return
	}

	render.JSON(w, r, user)
}

func (a *API) UpdateOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	updateDTO := &UpdateDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateDTO); err != nil {
		sendInternalServerError(w)
	}
	user, err := a.service.UpdateOne(id, updateDTO)
	if err != nil {
		//TODO
		return
	}
	render.JSON(w, r, user)
}

func (a *API) All(w http.ResponseWriter, r *http.Request) {
	dtos, err := a.service.All()
	if err != nil {
		sendInternalServerError(w)
		return
	}

	render.JSON(w, r, dtos)
}
