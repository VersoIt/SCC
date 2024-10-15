package model

import "time"

type EmployeeWithoutId struct {
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	HireDate  time.Time `db:"hire_date"`
}

type Employee struct {
	Id int `db:"id"`
	EmployeeWithoutId
}
