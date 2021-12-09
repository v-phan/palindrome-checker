package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/viptra/palindrom-ee/db"
	"github.com/viptra/palindrom-ee/util"
)

type PalindromeHandler struct {
	dbConn *db.DbConnection
}

type PalindromeResponse struct {
	FirstName bool `json: "firstName"`
	LastName  bool `json: "lastName"`
}

func NewPalindromeHandler(dbConn *db.DbConnection) *PalindromeHandler {
	return &PalindromeHandler{dbConn: dbConn}
}

func (p *PalindromeResponse) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(p)
}

func (h *PalindromeHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.checkUserForPalindrome(rw, r)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PalindromeHandler) checkUserForPalindrome(rw http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(r.Header.Get("UserId"), 10, 32)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to get userId from header: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.dbConn.GetUser(int32(userId))
	if err != nil {
		http.Error(rw, fmt.Sprintf("Unable to look up user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	resp := &PalindromeResponse{FirstName: util.CheckPalindrome(user.FirstName), LastName: util.CheckPalindrome(user.LastName)}
	err = resp.ToJSON(rw)
	if err != nil {
		log.Printf("error writing palindromeresponse to writer: %s", err)
	}
}
