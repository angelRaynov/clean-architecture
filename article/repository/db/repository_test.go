package db

import (
	"context"
	"fmt"
	"github.com/angelRaynov/clean-architecture/article/repository"
	"github.com/angelRaynov/clean-architecture/domain"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
	"time"
)
func TestArticleRepository_Fetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockArticles := []domain.Article{
		{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: domain.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: domain.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(mockArticles[0].ID, mockArticles[0].Title, mockArticles[0].Content,
			mockArticles[0].Author.ID, mockArticles[0].UpdatedAt, mockArticles[0].CreatedAt).
		AddRow(mockArticles[1].ID, mockArticles[1].Title, mockArticles[1].Content,
			mockArticles[1].Author.ID, mockArticles[1].UpdatedAt, mockArticles[1].CreatedAt)

	query := "SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE created_at > \\? ORDER BY created_at LIMIT \\?"
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := NewArticleRepository(db)
	cursor := repository.EncodeCursor(mockArticles[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestArticleRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id, title, content, author_id, updated_at, created_at FROM article WHERE id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := NewArticleRepository(db)

	num := int64(5)
	res, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestArticleRepository_Store(t *testing.T) {
	now := time.Now()
	ar := &domain.Article{
		Title: "Test",
		Content: "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.Author{
			ID: 1,
			Name: "Tolkien",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//TODO: fix the test
	query := "INSERT article SET title=\\?, content=\\?, author_id=\\?, updated_at=\\?, created_at=\\?"
	prep := mock.ExpectPrepare(query)
	fmt.Println(prep)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.CreatedAt).WillReturnResult(sqlmock.NewResult(12, 1))
	a := NewArticleRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)

}