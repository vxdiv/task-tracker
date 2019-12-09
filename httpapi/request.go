package httpapi

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func parseRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

const defaultPageSize = 20

type PagerRequestQuery struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

func BindPager(c echo.Context) (*PagerRequestQuery, error) {
	pager := &PagerRequestQuery{}
	if err := c.Bind(pager); err != nil {
		return nil, err
	}

	if pager.PageSize <= 0 {
		pager.PageSize = defaultPageSize
	}

	if pager.Page <= 0 {
		pager.Page = 1
	}

	return pager, nil
}

func (p PagerRequestQuery) Count() int {
	return p.PageSize
}

func (p PagerRequestQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type TimeRequestQuery struct {
	From string `query:"from"`
	To   string `query:"to"`
}

type TimeFilter struct {
	from time.Time
	to   time.Time
}

func BindTimeFilter(c echo.Context) (*TimeFilter, error) {
	query := &TimeRequestQuery{}
	if err := c.Bind(query); err != nil {
		return nil, err
	}

	filter := &TimeFilter{}
	if err := filter.parse(query.From, query.To); err != nil {
		return nil, err
	}

	return filter, nil
}

func (t TimeFilter) From() time.Time {
	return t.from
}

func (t TimeFilter) To() time.Time {
	return t.to
}

func (t *TimeFilter) parse(from, to string) (err error) {
	if t.from, err = parseUnixTime(from); err != nil {
		return err
	}

	if t.to, err = parseUnixTime(to); err != nil {
		return err
	}

	y, m, d := time.Now().Date()
	if t.to.Equal(time.Unix(0, 0)) {
		t.to = time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC)
	}

	return nil
}

func parseUnixTime(val string) (time.Time, error) {
	if val == "" {
		val = "0"
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(i, 0).UTC(), nil
}
