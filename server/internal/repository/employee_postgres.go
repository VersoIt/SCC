package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"serverClientClient/server/internal/model"
	"serverClientClient/server/pkg/database"
	"strconv"
	"time"
)

type EmployeePostgres struct {
	db *sqlx.DB
}

func NewEmployeePostgres(db *sqlx.DB) *EmployeePostgres {
	return &EmployeePostgres{db: db}
}

func (r *EmployeePostgres) GetById(id int) (model.Employee, error) {
	query := fmt.Sprintf(`SELECT * FROM %s e WHERE e.id = $1`, EmployeesTableName)
	var employee model.Employee
	err := r.db.Get(&employee, query, id)
	return employee, err
}

func (r *EmployeePostgres) Init(count int) (bool, error) {
	dbInitializer := database.NewInitializer[model.EmployeeWithoutId](r.db)

	initialized, err := dbInitializer.Init(EmployeesTableName, count, func(i int) model.EmployeeWithoutId {
		return model.EmployeeWithoutId{Email: fmt.Sprintf("email%d@gmail.com", i), FirstName: strconv.Itoa(i), LastName: strconv.Itoa(i), HireDate: time.Now()}
	})

	if err != nil {
		return false, err
	}

	return initialized, nil
}
