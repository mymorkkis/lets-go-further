package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, errorMessage string) {
	data := map[string]string{"error": errorMessage}
	app.serveJSON(w, r, status, data, nil)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	status := http.StatusInternalServerError
	app.errorResponse(w, r, status, http.StatusText(status))
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	status := http.StatusNotFound
	app.errorResponse(w, r, status, http.StatusText(status))
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	status := http.StatusMethodNotAllowed
	app.errorResponse(w, r, status, http.StatusText(status))
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
