package postgres

import (
	. "Contest/internal/domain"
	"database/sql"
	"errors"
	"fmt"
)

type TestRepository struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) AddItem(item Test) error {
	_, err := r.db.Exec("INSERT INTO tests(id, task_id, input, expected_result, points) VALUES ($1, $2, $3, $4, $5)",
		item.ID, item.TaskID, item.Input, item.ExpectedResult, item.Points)
	if err != nil {
		err = fmt.Errorf("In TestRepository(AddItem): %w", err)
	}
	return err
}

func (r *TestRepository) DeleteItem(id int) error {
	_, err := r.db.Exec("DELETE from tests where id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		err = fmt.Errorf("In TestRepository(DeleteItem): %w", err)
	}
	return err
}

func (r *TestRepository) UpdateItem(id int, newItem Test) error {
	_, err := r.db.Exec("UPDATE tests SET id=$1,task_id=$2, input=$3, expected_result=$4, points=$5 WHERE id=$6",
		newItem.ID, newItem.TaskID, newItem.Input, newItem.ExpectedResult, newItem.Points, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		err = fmt.Errorf("In TestRepository(UpdateItem): %w", err)
	}
	return err
}

func (r *TestRepository) GetTable() ([]Test, error) {
	rows, err := r.db.Query("SELECT * FROM tests")
	if err != nil {
		return nil, fmt.Errorf("In TestRepository(GetTable): %w", err)
	}
	defer rows.Close()

	tests := make([]Test, 0) //??
	for rows.Next() {
		var test Test
		err = rows.Scan(&test.ID, &test.TaskID, &test.Input, &test.ExpectedResult, &test.Points)
		if err != nil {
			return nil, fmt.Errorf("In TestRepository(GetTable): %w", err)
		}
		tests = append(tests, test)
	}
	return tests, nil
}

func (r *TestRepository) FindItemByID(id int) (Test, error) {
	row := r.db.QueryRow("SELECT * FROM tests WHERE id = $1", id)
	if row.Err() != nil {
		return Test{}, fmt.Errorf("In TestRepository(FindItemByID): %w", row.Err())
	}
	var test Test
	err := row.Scan(&test.ID, &test.TaskID, &test.Input, &test.ExpectedResult, &test.Points)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}
		err = fmt.Errorf("In TestRepository(FindItemByID): %w", err)
	}
	return test, err
}

func (r *TestRepository) FindItemByCondition(condition func(item Test) bool) (Test, error) {
	items, err := r.FindItemsByCondition(condition)
	if err != nil {
		return Test{}, fmt.Errorf("In TestRepository(FindItemByCondition): %w", err)
	}
	if len(items) == 0 {
		return Test{}, ErrNotFound
	}
	return items[0], nil
}

func (r *TestRepository) FindItemsByCondition(condition func(item Test) bool) ([]Test, error) {
	table, err := r.GetTable()
	if err != nil {
		return nil, fmt.Errorf("In TestRepository(FindItemsByCondition): %w", err)
	}
	res := make([]Test, 0) //???
	for _, test := range table {
		if condition(test) {
			res = append(res, test)
		}
	}
	return res, nil
}
