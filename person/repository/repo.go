package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test-task/person/entity"
	"time"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
}

type repository struct {
	db DBTX
}

type Repository interface {
	AddPerson(person *entity.Person) error
	GetPeople(nation string, page int, pageSize int) ([]*entity.Person, error)
	DeleteById(id int) error
	GetPersonById(id int) (entity.Person, error)
	UpdatePerson(id int, updated entity.Person) error
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (rp *repository) AddPerson(person *entity.Person) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var lastId int
	query := `insert into persons(name, surname, patronymic, age, gender, nationality) values ($1, $2, $3, $4, $5, $6) returning id`

	err := rp.db.QueryRowContext(ctx, query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&lastId)
	if err != nil {
		return fmt.Errorf("error in query executing: %v", err)
	}

	person.ID = lastId

	return nil
}

func (rp *repository) GetPeople(nation string, page int, pageSize int) ([]*entity.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `SELECT id, name, surname, patronymic,age, gender, nationality FROM persons where true`
	if nation != "" {
		query += fmt.Sprintf(" AND nationality = '%s'", nation)
	}
	offset := (page - 1) * pageSize

	query += fmt.Sprintf(" ORDER BY id LIMIT %d OFFSET %d", pageSize, offset)

	rows, err := rp.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying people: %s", err)
	}
	defer rows.Close()

	var people []*entity.Person

	for rows.Next() {
		var person entity.Person
		if err := rows.Scan(&person.ID, &person.Name,
			&person.Surname, &person.Patronymic, &person.Age,
			&person.Gender, &person.Nationality); err != nil {
			log.Println("Error scanning person:", err)
			return nil, err
		}
		people = append(people, &person)
	}

	return people, nil
}

func (rp *repository) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if !rp.personExists(id) {
		return fmt.Errorf("person with ID %d does not exist", id)
	}

	query := `DELETE FROM persons where id = $1`

	_, err := rp.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error in exec query: %s", err)
	}
	return nil
}

func (rp *repository) GetPersonById(id int) (entity.Person, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	person := entity.Person{}

	query := `Select id, name, surname, patronymic, age, gender, nationality from persons where id = $1`
	if err := rp.db.QueryRowContext(ctx, query, id).Scan(&person.ID, &person.Name, &person.Surname,
		&person.Patronymic, &person.Age, &person.Gender, &person.Nationality); err != nil {
		return entity.Person{}, fmt.Errorf("error scanning person: %s", err)
	}

	return person, nil
}

func (rp *repository) UpdatePerson(id int, updated entity.Person) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if !rp.personExists(id) {
		return fmt.Errorf("person with ID %d does not exist", id)
	}

	query := `UPDATE persons SET name=$2, surname=$3, patronymic=$4, age=$5, gender=$6, nationality=$7  WHERE id=$1`
	_, err := rp.db.ExecContext(ctx, query, id, updated.Name, updated.Surname, updated.Patronymic, updated.Age, updated.Gender, updated.Nationality)
	if err != nil {
		return fmt.Errorf("failed at query exec: %v", err)
	}

	return nil
}

func (rp *repository) personExists(id int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var count int
	query := "SELECT COUNT(*) FROM persons WHERE id = $1"
	err := rp.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		log.Println("Error checking if person exists:", err)
		return false
	}

	return count > 0
}
