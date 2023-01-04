package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.serveJSON(w, r, http.StatusOK, nil, nil)
}
