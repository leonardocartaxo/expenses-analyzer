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

func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	createDTO := &CreateDTO{}
	if err := json.NewDecoder(r.Body).Decode(createDTO); err != nil {
		//TODO
		//a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		//e.BadRequest(w, e.RespJSONDecodeFailure)
		return
	}
	newUser, err := a.service.Save(createDTO)
	if err != nil {
		//TODO
		return
	}
	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		//TODO
		return
	}
}

func (a *API) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := a.service.FindOne(id)
	if err != nil {
		//TODO
		return
	}
	render.JSON(w, r, user)
}

func (a *API) UpdateOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	updateDTO := &UpdateDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateDTO); err != nil {
		//TODO
		return
	}
	user, err := a.service.UpdateOne(id, updateDTO)
	if err != nil {
		//TODO
		return
	}
	render.JSON(w, r, user)
}

func (a *API) All(w http.ResponseWriter, r *http.Request) {
	models, err := a.service.All()
	if err != nil {
		//TODO
		return
	}
	var users []*DTO
	for _, user := range models {
		users = append(users, user.ToDTO())
	}

	//users2 := slices.AppendSeq(users, slices.Values(models))
	render.JSON(w, r, users)
}
