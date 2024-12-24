package types

// we make student's struct

type Student struct {
	Id    int
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"required"`
}