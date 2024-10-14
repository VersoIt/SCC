package db

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type Initializer[T any] struct {
	db *sqlx.DB
}

func (in *Initializer[T]) Init(table string, valueSets string, count int, applier func(int) T) error {
	if count <= 0 {
		return errors.New("count must be greater than 0")
	}

	isInit, err := in.isInitialized(table)
	if err != nil {
		return err
	}

	if isInit {
		return nil
	}

	firstElem := applier(0)
	elems := make([]T, 0, count)
	elems = append(elems, firstElem)
	for i := 1; i < count; i++ {
		elems = append(elems, applier(i))
	}

	elemType := reflect.TypeOf(firstElem)
	columnsToInsert := in.formColumnsToInsert(elemType)

}

func (in *Initializer[T]) isInitialized(table string) (bool, error) {
	count := 0
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if err := in.db.QueryRow(query).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (in *Initializer[T]) formColumnsToInsert(elemType reflect.Type) string {
	numOfFields := elemType.NumField()
	var columns []string

	for i := 0; i < numOfFields; i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			columns = append(columns, dbTag)
		} else {
			columns = append(columns, field.Name)
		}
	}

	return strings.Join(columns, ",")
}

func (in *Initializer[T]) formPlaceholdersToInsert(elemsLen, fieldsLen int) string {
	placeholders := make([]string, 0, elemsLen)
	for i := 0; i < elemsLen; i++ {
		fields := make([]string, 0, fieldsLen)
		for j := 1; j <= fieldsLen; j++ {
			fields = append(fields, fmt.Sprintf("$%d", j))
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(fields, ",")))
	}
	return strings.Join(placeholders, ",")
}

func (in *Initializer[T]) formArgsToInsert(elems []T) ([]interface{}, error) {
	if len(elems) == 0 {
		return nil, errors.New("elems is empty")
	}

	elemType := reflect.TypeOf(elems[0])

	args := make([]interface{}, 0, len(elems))
	for i := 0; i < len(elems); i++ {
		elem := elems[i]
		for j := 0; j < elemType.NumField(); j++ {

		}
	}
}
