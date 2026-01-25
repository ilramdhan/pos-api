package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination holds pagination parameters
type Pagination struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Sort    string `json:"sort"`
	Order   string `json:"order"`
}

// DefaultPagination contains default pagination values
var DefaultPagination = Pagination{
	Page:    1,
	PerPage: 10,
	Sort:    "created_at",
	Order:   "desc",
}

// MaxPerPage is the maximum allowed items per page
const MaxPerPage = 100

// GetPagination extracts pagination parameters from the request
func GetPagination(c *gin.Context) Pagination {
	p := DefaultPagination

	// Parse page
	if page := c.Query("page"); page != "" {
		if pageInt, err := strconv.Atoi(page); err == nil && pageInt > 0 {
			p.Page = pageInt
		}
	}

	// Parse per_page
	if perPage := c.Query("per_page"); perPage != "" {
		if perPageInt, err := strconv.Atoi(perPage); err == nil && perPageInt > 0 {
			p.PerPage = perPageInt
			if p.PerPage > MaxPerPage {
				p.PerPage = MaxPerPage
			}
		}
	}

	// Parse sort
	if sort := c.Query("sort"); sort != "" {
		p.Sort = sort
	}

	// Parse order
	if order := c.Query("order"); order != "" {
		if order == "asc" || order == "desc" {
			p.Order = order
		}
	}

	return p
}

// Offset returns the offset for SQL queries
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Limit returns the limit for SQL queries
func (p *Pagination) Limit() int {
	return p.PerPage
}

// OrderBy returns the ORDER BY clause for SQL queries
func (p *Pagination) OrderBy() string {
	return p.Sort + " " + p.Order
}
