package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/angelRaynov/clean-architecture/article/repository"
	"github.com/angelRaynov/clean-architecture/domain"
	"github.com/labstack/gommon/log"
)

type articleRepository struct {
	DB *sql.DB
}

func (ar *articleRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Article, err error) {
	rows, err := ar.DB.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Error(errRow)
		}
	}()

	result := make([]domain.Article, 0)
	for rows.Next() {
		var t domain.Article
		var authorID int64
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		t.Author = domain.Author{
			ID: authorID,
		}

		result = append(result, t)
	}

	return result, nil
}

func (ar *articleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at 
			FROM article WHERE created_at > ? ORDER BY created_at LIMIT ?`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadInput
	}

	res, err = ar.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, err
}

func (ar *articleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at 
			FROM article WHERE id = ?`

	list, err := ar.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	var res domain.Article

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return res, err
}

func (ar *articleRepository) GetByTitle(ctx context.Context, title string) (domain.Article, error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at 
			FROM article WHERE title = ?`

	list, err := ar.fetch(ctx, query, title)
	if err != nil {
		return domain.Article{}, err
	}

	var res domain.Article

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return res, err
}

func (ar *articleRepository) Update(ctx context.Context, a *domain.Article) error {
	query := `UPDATE article SET title=?, content=?, author_id=? updated_at=? WHERE id = ?`

	stmt, err := ar.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("err: rows affected %d", rowsAffected)
	}

	return err

}

func (ar *articleRepository) Store(ctx context.Context, a *domain.Article) error {
	query := `INSERT article SET title=?, content=?, updated_at=?, created_at=?`

	stmt, err := ar.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastID

	return err
}

func (ar *articleRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM article WHERE id = ?`

	stmt, err := ar.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("err: rows affected %d", rowsAffected)
	}

	return err
}

func newArticleRepository(db *sql.DB) domain.ArticleRepository {
	return &articleRepository{
		DB: db,
	}
}
