package engine

import (
    "fmt"
)

func PrintSwagger(schemas []Schema) {
    fmt.Println("📚 OpenAPI Spec (Generated Routes):")
    for _, schema := range schemas {
        prefix := "/api/" + schema.Model
        if schema.Routes.List {
            fmt.Printf("GET    %s\n", prefix)
        }
        if schema.Routes.Get {
            fmt.Printf("GET    %s/:id\n", prefix)
        }
        if schema.Routes.Create {
            fmt.Printf("POST   %s\n", prefix)
        }
        if schema.Routes.Update {
            fmt.Printf("PATCH  %s/:id\n", prefix)
        }
        if schema.Routes.Delete {
            fmt.Printf("DELETE %s/:id\n", prefix)
        }
    }
}
