package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"requestMethod": r.Method,
		"requestURL":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	data := map[string]any{"error": message}
	app.serveJSON(w, r, status, data, nil)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
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

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
