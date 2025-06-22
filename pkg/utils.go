package pkg

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func ParsePagination(c *gin.Context) (page, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ = strconv.Atoi(pageStr)
	pageSize, _ = strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return page, pageSize
}

func ParseFilters(c *gin.Context) map[string]string {
	filterStr := c.Query("filter")
	filters := make(map[string]string)

	if filterStr != "" {
		pairs := strings.Split(filterStr, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) == 2 {
				filters[kv[0]] = kv[1]
			}
		}
	}

	return filters
}

func ApplyPagination(db *gorm.DB, offset, pageSize int) *gorm.DB {
	return db.Offset(offset).Limit(pageSize)
}

func ApplyFilters(db *gorm.DB, filters map[string]string) *gorm.DB {
	for k, v := range filters {
		db = db.Where(k+" = ?", v)
	}
	return db
}
