// This file contains the repository implementation layer.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/commons"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

type Repository struct {
	Db *sql.DB
}

func (r *Repository) CreateUser(ctx context.Context, input UserInput) (int, error) {
	const queryTemplate = `
		INSERT INTO %s (phoneNumber, fullName, password, saltKey)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	query := fmt.Sprintf(queryTemplate, UserModel{}.TableName)

	var userID int
	if err := r.Db.QueryRowContext(ctx, query, input.PhoneNumber, input.FullName, input.Password, input.SaltKey).Scan(&userID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, errors.New(commons.ErrorNoData)
		default:
			return 0, err
		}
	}

	return userID, nil
}

func (r *Repository) GetUser(ctx context.Context, input GetUserInput) (*UserModel, error) {
	model := &UserModel{}

	where, args := BuildQuery(input)

	query := `
        SELECT
            id,
            phoneNumber,
            fullName,
            password,
            saltKey,
            createdAt,
            updatedAt
        FROM %s %s`

	query = fmt.Sprintf(query, model.TableName(), where)

	err := r.Db.QueryRowContext(ctx, query, args...).Scan(
		&model.ID,
		&model.PhoneNumber,
		&model.FullName,
		&model.Password,
		&model.SaltKey,
		&model.CreatedAt,
		&model.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(commons.ErrorNoData)
		}
		return nil, err
	}

	return model, nil
}

func (r *Repository) UpdateUser(ctx context.Context, input UserInput) error {
	//TODO implement me
	panic("implement me")
}

func BuildQuery(input interface{}) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()

	if inputType.Kind() != reflect.Struct {
		return "", args
	}

	argIndex := 1
	for i := 0; i < inputValue.NumField(); i++ {
		fieldValue := inputValue.Field(i)
		fieldType := inputType.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		if isNonNullPointer(fieldValue) {
			condition := fmt.Sprintf("%s = $%d", jsonTag, argIndex)
			conditions = append(conditions, condition)
			args = append(args, fieldValue.Elem().Interface())
			argIndex++
		}
	}

	if len(conditions) == 0 {
		return "", args
	}

	query := " WHERE " + strings.Join(conditions, " AND ")
	return query, args
}

// isNonNullPointer checks if a reflect.Value is a non-nil pointer.
func isNonNullPointer(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && !value.IsNil()
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		panic(err)
	}
	return &Repository{
		Db: db,
	}
}
