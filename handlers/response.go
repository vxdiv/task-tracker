package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

type ResponseItemList struct {
	CurrentPage int         `json:"current_page"`
	PagesCount  int         `json:"pages_count"`
	PageSize    int         `json:"page_size"`
	TotalCount  int         `json:"count"`
	Items       interface{} `json:"items"`
}

func ResponseItems(totalCount int, pager *Pager, items interface{}) ResponseItemList {
	return ResponseItemList{
		CurrentPage: pager.Page,
		TotalCount:  totalCount,
		PagesCount:  totalCount / pager.PageSize,
		PageSize:    pager.PageSize,
		Items:       items,
	}
}

func BadRequestError(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func NotFoundError(err error) error {
	return echo.NewHTTPError(http.StatusNotFound, err.Error())
}

func InternalServerError(err error) error {
	log.Errorf("Internal Server Error: %v", err)

	return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
}
