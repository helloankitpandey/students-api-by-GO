package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/helloankitpandey/students-api/internal/types"
	"github.com/helloankitpandey/students-api/internal/utils/response"
)

func New() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a student")

		// getting element of body in go
		// we don't directly get data but for getting we decode it first
		// for getting data in go we first serialize data in struct
		// so we need struct first for body -> define struct in types/types.go

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) // decode the sudent by taking student's address 
		// and store it in err for it returning err
		if errors.Is(err, io.EOF) { // EOF is the error returned by Read when no more input is available.
			// agr no input hai to
			//  we return response in josn formate 
			// so for json response make package => utils/response/response.go

			// response.WriteJson(w, http.StatusBadRequest, err.Error())
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err)) // or we return as "empty-body" 
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty -body")))
			return 
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// Now we do request validation

		// we do it manualy but go has paowerful package for it
		// use this go get github.com/go-playground/validator/v10
		// then in types/types.go write `validate:"required"` after string written

		if err := validator.New().Struct(student); err != nil {
			
			// typcast the types
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}


		

		// w.Write([]byte("welcome to students api"))
		// here we return response as json
		response.WriteJson(w, http.StatusCreated, map[string]string{"Success": "Ok"})
	}
}