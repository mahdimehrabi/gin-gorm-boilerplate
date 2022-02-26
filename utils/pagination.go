package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

//Pagination -> struct for Pagination
type Pagination struct {
	Page       int
	Sort       string
	PageSize   int
	Offset     int
	All        bool
	Keyword    string
}

//BuildPagination -> builds the pagination
func BuildPagination(c *gin.Context) Pagination {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")
	sort := c.Query("sort")
	keyword := c.Query("keyword")

	var all bool
	if pageSizeStr == "Infinity" {
		all = true
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	return Pagination{
		Page:       page,
		Sort:       sort,
		PageSize:   pageSize,
		Offset:     (page - 1) * pageSize,
		All:        all,
		Keyword:    keyword,
	}
}