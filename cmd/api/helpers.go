package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SystemInfo struct {
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type JsonResponse struct {
	Status     Status     `json:"status"`
	SystemInfo SystemInfo `json:"systemInfo"`
	Data       any        `json:"data"`
}

func (app *application) serveJSON(w http.ResponseWriter, r *http.Request, status int, data any, headers http.Header) {
	json_, err := json.Marshal(&JsonResponse{
		Status: Status{
			Code:    status,
			Message: http.StatusText(status),
		},
		SystemInfo: SystemInfo{
			Environment: app.config.env,
			Version:     app.version,
		},
		Data: data,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json_)
}

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int) {
	app.serveJSON(w, r, status, nil, nil)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	app.errorResponse(w, r, http.StatusInternalServerError)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed)
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
