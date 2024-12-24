package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// define type of Response for using error returning
type Response struct {
	Status string `json:"status"` // serializing in json as status not Status
	Error  string `json:"error"`
}

// make constant for returning 
const (
	StatusOK = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface {}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// /before it we decode but here we encoded
	return json.NewEncoder(w).Encode(data) // encode method return error type
}


// error ko josn formate me return krne ke liye 
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}

// validation error ke liye ak function
// ye errs ke package ki list send krta hai

func ValidationError(errs validator.ValidationErrors) Response {
	
	var errMsgs []string

	// loop
	for _,err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field())) //sprintf used for concatinat string
			// field() -> give us kon se validate require hai like name,email or age
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs, ","),
	}
}