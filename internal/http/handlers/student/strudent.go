package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	// "github.com/helloankitpandey/students-api/internal/http/handlers/student"
	"github.com/helloankitpandey/students-api/internal/storage"
	"github.com/helloankitpandey/students-api/internal/types"
	"github.com/helloankitpandey/students-api/internal/utils/response"
)

// route for posting student data to database
func New(storage storage.Storage) http.HandlerFunc {

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

		// Now Student ko create karenge using storage package
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User Created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}


		// w.Write([]byte("welcome to students api"))
		// here we return response as json
		// response.WriteJson(w, http.StatusCreated, map[string]string{"Success": "Ok"})
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})

	}
}

// Now new route for getting data of student by id
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		// Now we need a method over storage for getting a students
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		
		//handle database error 
		if err != nil {
			// add this error logging in everywhere where its needed
			slog.Error("error getting user", slog.String("id", id)) // add this error logging in everywhere where its needed
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

// Now new routw for getting all students data in one click
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		// Now make interface/method in storage.go file
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)
		
		
	}
}
fmt.Println("cfg is loaded")
