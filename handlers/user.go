package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/viptra/palindrom-ee/db"
)

type UserHandler struct {
	dbConn *db.DbConnection
}

type UserRequest struct {
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
}

func NewUserHandler(d *db.DbConnection) *UserHandler {
	return &UserHandler{dbConn: d}
}

func (u *UserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		u.getUser(rw, r)

	} else if r.Method == http.MethodPost {
		u.createUser(rw, r)

	} else if r.Method == http.MethodPut {
		u.updateUser(rw, r)

	} else if r.Method == http.MethodDelete {
		u.deleteUser(rw, r)

	} else {
		// Siden bruker kan bruke flere operasjoner
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (u *UserHandler) createUser(rw http.ResponseWriter, r *http.Request) {
	newUser := &UserRequest{}
	//Unmarshal requestbody to userreq. struct
	err := json.NewDecoder(r.Body).Decode(newUser)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to read user details in your request: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if newUser.FirstName == "" || newUser.LastName == "" {
		http.Error(rw, "Firstname or lastname was not defined in your payload.", http.StatusBadRequest)
		return
	}
	log.Println("Creating user")
	user, err := u.dbConn.CreateUser(newUser.FirstName, newUser.LastName)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	user.ToJSON(rw)
}

func (u *UserHandler) updateUser(rw http.ResponseWriter, r *http.Request) {
	newUser := &UserRequest{}

	userId, err := strconv.ParseInt(r.Header.Get("UserId"), 10, 32)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//Unmarshal requestbody to userreq. struct
	err = json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//Users have to send in both names
	if newUser.FirstName == "" || newUser.LastName == "" {
		http.Error(rw, "Firstname or lastname was not defined in your payload.", http.StatusBadRequest)
		return
	}

	_, err = u.dbConn.GetUser(int32(userId))
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Creating user")
			user, err := u.dbConn.CreateUser(newUser.FirstName, newUser.LastName)
			if err != nil {
				http.Error(rw, fmt.Sprintf("Unable to create user: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			user.ToJSON(rw)
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := u.dbConn.UpdateUser(int32(userId), newUser.FirstName, newUser.LastName)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to update user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	user.ToJSON(rw)
}
func (u *UserHandler) deleteUser(rw http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(r.Header.Get("UserId"), 10, 32)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	deleted, err := u.dbConn.DeleteUser(int32(userId))
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to delete user %v, error: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(rw, fmt.Sprintf("Unable to delete user %v, user does not exist", userId), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (u *UserHandler) getUser(rw http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(r.Header.Get("UserId"), 10, 32)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := u.dbConn.GetUser(int32(userId))
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(rw, "Unable to find user", http.StatusNoContent)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	user.ToJSON(rw)
}
