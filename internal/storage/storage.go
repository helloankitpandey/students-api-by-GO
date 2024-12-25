package storage

// interface ka use krke ham apne application ko plugin type key bna skte 
// means if we uses sql-lite and we want to switch in postgres then it will done easily by minimal changes && also good for testing
type Storage interface{
	CreateStudent(name string, email string, age int) (int64, error)
}