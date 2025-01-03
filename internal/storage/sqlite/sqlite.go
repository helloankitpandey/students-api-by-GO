package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/helloankitpandey/students-api/internal/config"
	"github.com/helloankitpandey/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

// here we have to implement Storage interface

type Sqlite struct {
	Db *sql.DB
}

// splite ka instatnce banane ke liye ham Constructor banate hai in oops
// but here we make func
func New(cfg *config.Config) (*Sqlite, error) {

	db,err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

// At 2:06:24 	we got error here => router.HandleFunc("POST /api/students", student.New(storage))
// So for solving this we implement storage interface

// Now our Sqlite s struct implement the Storage interface
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	
	// prepare sql but not passes values for sql injection se bachne ke liye 
	stmt, err := s.Db.Prepare("INSERT INTO students(name, email, age) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close() // prepare ko excute hone ke bad close krna hota hai

	result, err := stmt.Exec(name, email, age) // storing data in squential way
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId() // now use and return id by using LastInserId method
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

// now after makinf New we have to call it main.go file


// Now implenting new interface 
// i.e etStudentById(id int64) (types.Student, error)
// for getting student data by id
func (s *Sqlite) GetStudentById(id int64) ( types.Student , error) {

	// sbse phle sql query ko prepare krna hoga
	stmt, err := s.Db.Prepare("SELECT id, name, age, email FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err 
	}

	defer stmt.Close()

	// Now jo data database se aa rha hai usse hame struct ke andar serialize krke dalna hoga
	var student types.Student

	// same as table created in sql same column-wise
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Age, &student.Email) // same as table created in sql same column-wise

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

	// now getstudentbyid is ready to use in student.go/GetById 
}


// NOw implementing new interface
// i.e 	GetStudents() ([]types.Student, error)
// for getting all students data
func (s *Sqlite) GetStudents() ([]types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, age, email FROM students")
	if err != nil {
		return nil ,err
	}

	defer stmt.Close()

	// now execute it
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Age, &student.Email)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	// returning all students from looping
	return students, nil
}

