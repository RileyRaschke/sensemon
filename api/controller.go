package api

import (
	"net/http"
	"strconv"

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
	r.Get("/sensors", c.sensors)
	r.Get("/sensors/latest", c.latestSensorReadings)
	r.Get("/sensordata/{sensor_id}", c.sensorData)
	r.Get("/sensordata/{sensor_id}/{interval}", c.sensorDataInterval)
	r.Get("/allsensors/{interval}", c.allSensorsInterval)
	c.router = r
	return r
}

func (c *ApiController) sensors(w http.ResponseWriter, r *http.Request) {
	data, err := c.dbc.Sensors()
	if err != nil {
		log.Errorf("DB Error: %w", err)
		view.JsonErrorMsg(w, http.StatusInternalServerError, "Database Error")
		return
	}
	view.AsJson(w, data)
}

func (c *ApiController) latestSensorReadings(w http.ResponseWriter, r *http.Request) {
	data, err := c.dbc.LatestDhtReadings()
	if err != nil {
		log.Errorf("DB Error: %w", err)
		view.JsonErrorMsg(w, http.StatusInternalServerError, "Database Error")
		return
	}
	view.AsJson(w, data)
}

func (c *ApiController) sensorData(w http.ResponseWriter, r *http.Request) {
	sensor_id, err := url.QueryUnescape(chi.URLParam(r, "sensor_id"))
	if err != nil {
		log.Errorf("Bad Request: %w", err)
		view.JsonErrorMsg(w, http.StatusBadRequest, "Bad Request")
		return
	}
	log.Infof("Searching for sensor_id: %s", sensor_id)
	data, err := c.dbc.AllDhtDataForSensor(sensor_id)

	if err != nil {
		log.Errorf("DB Error: %w", err)
		view.JsonErrorMsg(w, http.StatusInternalServerError, "Database Error")
		return
	}
	view.AsJson(w, data)
}

func (c *ApiController) sensorDataInterval(w http.ResponseWriter, r *http.Request) {
	sensor_id, err := url.QueryUnescape(chi.URLParam(r, "sensor_id"))
	interval, err2 := strconv.Atoi(chi.URLParam(r, "interval"))
	if err != nil || err2 != nil {
		log.Errorf("Bad Request: %w", err)
		view.JsonErrorMsg(w, http.StatusBadRequest, "Bad Request")
		return
	}

	log.Infof("Searching for sensor_id: %s", sensor_id)
	data, err := c.dbc.AllDhtDataForSensorInterval(sensor_id, interval)

	if err != nil {
		log.Errorf("DB Error: %w", err)
		view.JsonErrorMsg(w, http.StatusInternalServerError, "Database Error")
		return
	}
	view.AsJson(w, data)
}

func (c *ApiController) allSensorsInterval(w http.ResponseWriter, r *http.Request) {
	interval, err := strconv.Atoi(chi.URLParam(r, "interval"))
	if err != nil {
		log.Errorf("Bad Request: %w", err)
		view.JsonErrorMsg(w, http.StatusBadRequest, "Bad Request")
		return
	}
	data, err := c.dbc.AllDhtDataInterval(interval)

	if err != nil {
		log.Errorf("DB Error: %w", err)
		view.JsonErrorMsg(w, http.StatusInternalServerError, "Database Error")
		return
	}
	view.AsJson(w, data)
}
