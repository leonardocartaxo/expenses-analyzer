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
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		//TODO
		//a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		//e.BadRequest(w, e.RespJSONDecodeFailure)
		return
	}
	newUser := a.service.Save(form.Name)
	err := json.NewEncoder(w).Encode(newUser)
	if err != nil {
		//TODO
		return
	}

}

func (a *API) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := a.service.FindOne(id)
	render.JSON(w, r, user)
}

func (a *API) All(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, a.service.All())
}
