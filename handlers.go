package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sa-course-app/pkg/user/app"
	"sa-course-app/pkg/user/domain"
)

type handlers struct {
	s *app.Service
}

func (h *handlers) readyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "{\"host\": \"%v\"}", r.Host)
}

func (h *handlers) healthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func (h *handlers) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userData app.UserData

	err := decoder.Decode(&userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
	} else {
		id, err := h.s.CreateUser(&userData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v", err)
		} else {
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintf(w, "%v", uuid.UUID(id).String())
		}
	}
}

func (h *handlers) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid id \"%s\"", id)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var userData app.UserData

	err = decoder.Decode(&userData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	err = h.s.UpdateUser(domain.UserID(uid), &userData)
	if err == app.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

func (h *handlers) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid id \"%s\"", id)
		return
	}

	userID := domain.UserID(uid)
	userData, err := h.s.FindUser(userID)
	if err != nil {
		if err == app.ErrUserNotFound {
			w.WriteHeader(http.StatusNotFound)
		}

		log.Printf("%v", err)
		return
	}

	d := json.NewEncoder(w)

	w.WriteHeader(http.StatusOK)
	err = d.Encode(userData)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (h *handlers) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Invalid id \"%s\"", id)
		return
	}

	err = h.s.DeleteUser(domain.UserID(uid))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
