package usecase

import (
	"context"
	"github.com/angelRaynov/clean-architecture/domain"
	"github.com/labstack/gommon/log"
	"time"
	"golang.org/x/sync/errgroup"
)

type articleUseCase struct {
	articleRepo domain.ArticleRepository
	authorRepo domain.AuthorRepository
	contextTimeout time.Duration
}

func NewArticleUseCase(a domain.ArticleRepository, ar domain.AuthorRepository, timeout time.Duration) domain.ArticleUseCase {
	return &articleUseCase{
		articleRepo: a,
		authorRepo: ar,
		contextTimeout: timeout,
	}
}

func (a articleUseCase) Fetch(ctx context.Context, cursor string, num int64) ([]domain.Article, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(ctx,a.contextTimeout)
	defer cancel()

	res, nextCursor, err := a.articleRepo.Fetch(ctx,cursor,num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx,res)
	if err != nil {
		nextCursor = ""
	}

	return res,nextCursor, err
}

func (a articleUseCase) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return res, nil
}

func (a articleUseCase) Update(ctx context.Context, ar *domain.Article) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)

	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a articleUseCase) GetByTitle(ctx context.Context, title string) (domain.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)

	defer cancel()

	res, err := a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return domain.Article{},err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Article{}, err
	}

	res.Author = resAuthor
	return res, nil
}

func (a articleUseCase) Store(ctx context.Context, article *domain.Article) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)

	defer cancel()

	existingArticle, _ := a.GetByTitle(ctx, article.Title)
	if existingArticle != (domain.Article{}) {
		return domain.ErrConflict
	}

	err := a.articleRepo.Store(ctx, article)
	return err
}

func (a articleUseCase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)

	defer cancel()

	existingArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingArticle == (domain.Article{}) {
		return domain.ErrNotFound
	}

	return a.articleRepo.Delete(ctx, id)
}

func (a *articleUseCase) fillAuthorDetails(c context.Context, data []domain.Article) ([]domain.Article, error) {
	g, ctx := errgroup.WithContext(c)

	mapAuthors := map[int64]domain.Author{}

	for _, article := range data {
		mapAuthors[article.Author.ID] = domain.Author{}
	}

	chanAuthor := make(chan domain.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}
	go func() {
		err := g.Wait()
		if err != nil {
			log.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (domain.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}

	return data,nil
}
