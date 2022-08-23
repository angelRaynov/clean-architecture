package db

import (
	"context"
	"database/sql"
	"github.com/angelRaynov/clean-architecture/domain"
)

type authorRepo struct {
	DB *sql.DB
}

func NewAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &authorRepo{
		DB: db,
	}
}

func (a *authorRepo) getOne(ctx context.Context, query string, args ...interface{}) (domain.Author, error) {
	stmt, err := a.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Author{}, err
	}

	row := stmt.QueryRowContext(ctx, args...)

	var res domain.Author

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	return res, err
}
func (a *authorRepo) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	query := `SELECT id,name,created_at,updated_at FROM author WHERE id=?`
	return a.getOne(ctx, query, id)
}
