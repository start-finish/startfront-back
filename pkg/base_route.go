package pkg

import (
	"github.com/gin-gonic/gin"
)

type RouteOptions struct {
	EnableList   bool
	EnableGet    bool
	EnableCreate bool
	EnableUpdate bool
	EnableDelete bool
}

func RegisterCRUDRoutes[T any](
	group *gin.RouterGroup,
	svc *BaseService[T],
	opts RouteOptions,
	uniqueFields []string,
) {
	if opts.EnableList {
		group.GET("", listHandler(svc))
	}
	if opts.EnableGet {
		group.GET("/:id", getByIDHandler(svc))
	}
	if opts.EnableCreate {
		group.POST("", createHandler(svc, uniqueFields))
	}
	if opts.EnableUpdate {
		group.PATCH("/:id", updateHandler(svc))
	}
	if opts.EnableDelete {
		group.DELETE("/:id", deleteHandler(svc))
	}
}

func createHandler[T any](svc *BaseService[T], uniqueFields []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item T
		if err := c.ShouldBindJSON(&item); err != nil {
			errorWithStatus(c, "invalid data")
			return
		}

		for _, field := range uniqueFields {
			exists, err := svc.Exists(field, item)
			if err != nil {
				errorWithStatus(c, "validation failed")
				return
			}
			if exists {
				errorWithStatus(c, "duplicate "+field)
				return
			}
		}

		if err := svc.Create(&item); err != nil {
			errorWithStatus(c, "failed to create")
			return
		}

		success(c, "created successfully", item)
	}
}

func listHandler[T any](svc *BaseService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := svc.List()
		if err != nil {
			errorWithStatus(c, "failed to list data")
			return
		}
		success(c, "", items)
	}
}

func getByIDHandler[T any](svc *BaseService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri struct {
			ID uint `uri:"id"`
		}
		if err := c.ShouldBindUri(&uri); err != nil {
			errorWithStatus(c, "invalid id")
			return
		}

		item, err := svc.GetByID(uri.ID)
		if err != nil {
			errorWithStatus(c, "data not found")
			return
		}

		success(c, "", item)
	}
}

func updateHandler[T any](svc *BaseService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri struct {
			ID uint `uri:"id"`
		}
		if err := c.ShouldBindUri(&uri); err != nil {
			errorWithStatus(c, "invalid id")
			return
		}

		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			errorWithStatus(c, "invalid update data")
			return
		}

		if err := svc.Update(uri.ID, data); err != nil {
			errorWithStatus(c, err.Error())
			return
		}

		success(c, "updated successfully", nil)
	}
}

func deleteHandler[T any](svc *BaseService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri struct {
			ID uint `uri:"id"`
		}
		if err := c.ShouldBindUri(&uri); err != nil {
			errorWithStatus(c, "invalid ID")
			return
		}

		if err := svc.Delete(uri.ID); err != nil {
			errorWithStatus(c, err.Error())
			return
		}

		success(c, "deleted successfully", nil)
	}
}
