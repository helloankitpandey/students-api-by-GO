package sqlite

import (
	"database/sql"

	"github.com/helloankitpandey/students-api/internal/config"
	_"github.com/mattn/go-sqlite3" 
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