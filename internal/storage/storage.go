package storage

import "github.com/helloankitpandey/students-api/internal/types"

// interface ka use krke ham apne application ko plugin type key bna skte
// means if we uses sql-lite and we want to switch in postgres then it will done easily by minimal changes && also good for testing
type Storage interface{
	fmt.Println("cfg is loaded")

	CreateStudent(name string, email string, age int) (int64, error)
	// inteface for get student data by id & implement it in sqlite.go by making a method of same name
	GetStudentById(id int64) (types.Student, error)
	// interface for getting all students data & don't take any input && implement in slite package ke ander
	GetStudents() ([]types.Student, error)
}