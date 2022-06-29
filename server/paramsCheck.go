package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func isBodyEmpty(w http.ResponseWriter, r *http.Request) bool {
	if r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Body Required")
		return true
	}
	return false
}

func isContentJson(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") == ContentTypeJSON {
		return true
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Content-Type must be \"%s\"", ContentTypeJSON)
	return false
}

func isContentURLENCODED(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") == ContentTypeURLENCODED {
		return true
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Content-Type must be \"%s\"", ContentTypeURLENCODED)
	return false
}

func isSameMethod(w http.ResponseWriter, r *http.Request, m string) bool {
	if r.Method == m {
		return true
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Method not allowed")
	return false
}

func urlEncReqHasValidSettings(w http.ResponseWriter, r *http.Request) bool {
	if !isSameMethod(w, r, "POST") || isBodyEmpty(w, r) || !isContentURLENCODED(w, r) {
		return false
	}
	r.ParseForm()
	if len(r.PostForm) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Multiple values detected, the only valid input is 1 \"ID=X\"")
		return false
	}
	if _, ok := r.PostForm["ID"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid key detected, the only valid input is 1 \"ID=X\"")
		return false
	}
	if len(r.PostForm["ID"]) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only one ID is allowed, and it must be an int")
		return false
	}
	val, err := strconv.Atoi(r.PostForm["ID"][0])
	if err != nil || val < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID must be a positive int")
		return false
	}
	return true
}
