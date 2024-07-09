package main

import "net/http"

func errorHandler(w http.ResponseWriter, r *http.Request) {

	respondError(w, 400, "ERROR! Something went wrong!")

}
