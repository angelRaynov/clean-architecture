package http

import (
	"github.com/angelRaynov/clean-architecture/domain"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	ArticleUseCase domain.ArticleUseCase
}

func NewArticleHandler(e *echo.Echo, useCase domain.ArticleUseCase) {
	handler := &ArticleHandler{
		ArticleUseCase: useCase,
	}

	e.GET("/articles", handler.FetchArticle)
	e.GET("/articles/:id", handler.GetByID)
	e.POST("/articles", handler.Store)
	e.DELETE("/articles/:id", handler.Delete)
}

func (ah *ArticleHandler) FetchArticle(ec echo.Context) error {
	numString := ec.QueryParam("num")
	num, _ := strconv.Atoi(numString)

	cursor := ec.QueryParam("cursor")
	ctx := ec.Request().Context()

	listArticle, nextCursor, err := ah.ArticleUseCase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return ec.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	ec.Response().Header().Set(`X-Cursor`, nextCursor)
	return ec.JSON(http.StatusOK, listArticle)
}

func (ah *ArticleHandler) GetByID(ec echo.Context) error {
	idString, err := strconv.Atoi(ec.Param("id"))

	if err != nil {
		return ec.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idString)
	ctx := ec.Request().Context()

	article, err := ah.ArticleUseCase.GetByID(ctx, id)
	if err != nil {
		return ec.JSON(getStatusCode(err), ResponseError{
			Message: err.Error(),
		})
	}

	return ec.JSON(http.StatusOK, article)
}

func (ah *ArticleHandler) Store(ec echo.Context) error {
	var article domain.Article
	err := ec.Bind(&article)
	if err != nil {
		return ec.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isValidRequest(&article); !ok {
		return ec.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := ec.Request().Context()
	err = ah.ArticleUseCase.Store(ctx, &article)
	if err != nil {
		return ec.JSON(getStatusCode(err), ResponseError{
			Message: err.Error(),
		})
	}

	return ec.JSON(http.StatusCreated, article)
}

func (ah *ArticleHandler) Delete(ec echo.Context) error {
	idString, err := strconv.Atoi(ec.Param("id"))
	if err != nil {
		return ec.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idString)
	ctx := ec.Request().Context()

	err = ah.ArticleUseCase.Delete(ctx, id)
	if err != nil {
		ec.JSON(getStatusCode(err), ResponseError{
			Message: err.Error(),
		})
	}

	return ec.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError

	}
}

func isValidRequest(a *domain.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(a)
	if err != nil {
		return false,err
	}

	return true,nil
}