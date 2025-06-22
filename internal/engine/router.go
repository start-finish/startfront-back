package engine

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, schema Schema) {
    prefix := "/api/" + schema.Model

    if schema.Routes.List {
        r.GET(prefix, DynamicListHandler(schema))
    }
    if schema.Routes.Get {
        r.GET(prefix+"/:id", DynamicGetHandler(schema))
    }
    if schema.Routes.Create {
        r.POST(prefix, DynamicCreateHandler(schema))
    }
    if schema.Routes.Update {
        r.PATCH(prefix+"/:id", DynamicUpdateHandler(schema))
    }
    if schema.Routes.Delete {
        r.DELETE(prefix+"/:id", DynamicDeleteHandler(schema))
    }
}
