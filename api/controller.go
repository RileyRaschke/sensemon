package api

import (
	"net/http"

	"net/url"

	"sensemon/db"
	"sensemon/view"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type ApiController struct {
	dbc    *db.Connection
	router *chi.Mux
}

func NewApiController(dbc *db.Connection) *ApiController {
	return &ApiController{dbc, nil}
}

func (c *ApiController) Handler() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/sensordata/{sensor_id}", c.sensorData)
	c.router = r
	return r
}

func (c *ApiController) sensorData(w http.ResponseWriter, r *http.Request) {
	sensor_id, err := url.QueryUnescape(chi.URLParam(r, "sensor_id"))
	if err != nil {
		log.Errorf("Bad Request: %w", err)
		view.JsonErrorMsg(w, http.StatusBadRequest, "Bad Request")
		return
	}
	log.Infof("Searching for sensor_id: %s", sensor_id)
	data, err := c.dbc.AllTables()

	if err != nil {
		log.Errorf("LDAP Error: %w", err)
		view.JsonErrorMsg(w, http.StatusInternalServerError, "Database Error")
		return
	}
	view.AsJson(w, data)
}
