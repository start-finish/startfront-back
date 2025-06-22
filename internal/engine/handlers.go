package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/start-finish/startfront-app/pkg"
	"github.com/start-finish/startfront-app/pkg/response"
)

func DynamicListHandler(schema Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := pkg.GetDB()
		page, pageSize := pkg.ParsePagination(c)
		filters := pkg.ParseFilters(c)

		var results []map[string]interface{}

		db = db.Table(schema.Model)
		db = pkg.ApplyFilters(db, filters)
		db = pkg.ApplyPagination(db, (page-1)*pageSize, pageSize)

		if err := db.Find(&results).Error; err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, "Query successful", results)
	}
}

func DynamicGetHandler(schema Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := pkg.GetDB()
		id := c.Param("id")

		var result map[string]interface{}
		if err := db.Table(schema.Model).First(&result, id).Error; err != nil {
			response.Error(c, "Record not found")
			return
		}

		response.Success(c, "Query successful", result)
	}
}

func DynamicCreateHandler(schema Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := pkg.GetDB()
		var data map[string]interface{}

		if err := c.BindJSON(&data); err != nil {
			response.Error(c, "Invalid data")
			return
		}

		model, err := BuildGormModel(schema)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		// Set fields on model
		SetModelFields(model, data)

		if err := db.Table(schema.Model).Create(&data).Error; err != nil {
			response.Error(c, "Save failed: "+err.Error())
			return
		}

		response.Success(c, "Created", data)
	}
}

func DynamicUpdateHandler(schema Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := pkg.GetDB()
		id := c.Param("id")
		var data map[string]interface{}

		if err := c.BindJSON(&data); err != nil {
			response.Error(c, "Invalid data")
			return
		}

		if err := db.Table(schema.Model).Where("id = ?", id).Updates(data).Error; err != nil {
			response.Error(c, "Update failed: "+err.Error())
			return
		}

		response.Success(c, "Updated", nil)
	}
}

func DynamicDeleteHandler(schema Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := pkg.GetDB()
		id := c.Param("id")

		if err := db.Table(schema.Model).Where("id = ?", id).Delete(nil).Error; err != nil {
			response.Error(c, "Delete failed: "+err.Error())
			return
		}

		response.Success(c, "Deleted", nil)
	}
}
