package database

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

func NewInitializer[T any](db *sqlx.DB) *Initializer[T] {
	return &Initializer[T]{db: db}
}

func (in *Initializer[T]) Init(table string, count int, applier func(int) T) (bool, error) {
	if count <= 0 {
		return false, errors.New("count must be greater than 0")
	}

	isInit, err := in.isInitialized(table)
	if err != nil {
		return false, err
	}

	if isInit {
		return false, nil
	}

	firstElem := applier(0)
	elems := make([]T, 0, count)
	elems = append(elems, firstElem)
	for i := 1; i < count; i++ {
		elems = append(elems, applier(i))
	}

	elemType := reflect.TypeOf(firstElem)

	columnsToInsertChan := make(chan string)
	go func() {
		columnsToInsertChan <- in.formColumnsToInsert(elemType)
	}()

	placeholdersToInsertChan := make(chan string)
	go func() {
		placeholdersToInsertChan <- in.formPlaceholdersToInsert(len(elems), in.getFieldsLenOfType(elemType))
	}()

	argsToInsertChan := make(chan []any)
	argsToInsertChanErr := make(chan error)
	go func(elems []T) {
		data, err := in.formArgsToInsert(elems)
		if err != nil {
			argsToInsertChanErr <- err
			return
		}
		argsToInsertChan <- data
	}(elems)

	columnsToInsert := <-columnsToInsertChan
	placeholdersToInsert := <-placeholdersToInsertChan
	var argsToInsert []any
	select {
	case err := <-argsToInsertChanErr:
		return false, err
	case argsToInsert = <-argsToInsertChan:
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, columnsToInsert, placeholdersToInsert)
	_, err = in.db.Exec(query, argsToInsert...)
	return true, err
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
		if dbTag == "-" {
			continue
		} else if dbTag != "" {
			columns = append(columns, dbTag)
		} else {
			columns = append(columns, field.Name)
		}
	}

	return strings.Join(columns, ",")
}

func (in *Initializer[T]) getFieldsLenOfType(t reflect.Type) int {
	count := 0
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("db") == "-" {
			continue
		}
		count++
	}
	return count
}

func (in *Initializer[T]) formPlaceholdersToInsert(elemsLen, fieldsLen int) string {
	placeholders := make([]string, 0, elemsLen)
	count := 1
	for i := 0; i < elemsLen; i++ {
		fields := make([]string, 0, fieldsLen)
		for j := 1; j <= fieldsLen; j++ {
			fields = append(fields, fmt.Sprintf("$%d", count))
			count++
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
			tag := elemType.Field(j).Tag.Get("db")
			if tag == "-" {
				continue
			}
			args = append(args, reflect.ValueOf(elem).Field(j).Interface())
		}
	}
	return args, nil
}
