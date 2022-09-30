package usecase

import (
	"context"
	"errors"
	"github.com/angelRaynov/clean-architecture/domain"
	"github.com/angelRaynov/clean-architecture/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestArticleUseCase_Fetch(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title: "Hello",
		Content: "Content",
	}

	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArticle, "next-cursor", nil).Once()
		mockAuthor := domain.Author{
			ID: 1,
			Name: "King",
		}

		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second * 2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArticle))

		mockAuthorRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)

		t.Run("error", func(t *testing.T) {
			mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
				mock.AnythingOfType("int64")).Return(nil, "", errors.New("unexpected error")).Once()
			mockAuthorRepo = new(mocks.AuthorRepository)
			u = NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second * 2)
			num = int64(1)
			cursor = "12"
			list, nextCursor, err = u.Fetch(context.TODO(), cursor, num)

			assert.Empty(t, nextCursor)
			assert.Error(t, err)
			assert.Len(t, list, 0)
			mockArticleRepo.AssertExpectations(t)
			mockAuthorRepo.AssertExpectations(t)
		})
	})
}

func TestArticleUseCase_GetByID(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title: "Hello",
		Content: "Content",
	}

	mockAuthor := domain.Author{
		ID: 1,
		Name: "Tolkien",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second * 2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, errors.New("unexpected err")).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second * 2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.Article{}, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestArticleUseCase_Store(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(domain.Article{}, domain.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Article")).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Martin",
		}
		mockAuthorRepo := new(mocks.AuthorRepository)
		mockAuthorRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})

}

func TestArticleUseCase_Delete(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()

		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
	t.Run("article-does-not-exist", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, nil).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
	t.Run("db-error", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, errors.New("Unexpected Error")).Once()

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorRepo.AssertExpectations(t)
	})
}

func TestArticleUseCase_Update(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Update", mock.Anything, &mockArticle).Once().Return(nil)

		mockAuthorRepo := new(mocks.AuthorRepository)
		u := NewArticleUseCase(mockArticleRepo, mockAuthorRepo, time.Second*2)

		err := u.Update(context.TODO(), &mockArticle)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}