package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadSchemas(path string) ([]Schema, error) {
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return nil, err
    }

    var schemas []Schema

    for _, file := range files {
        if !file.IsDir() && (file.Mode().IsRegular()) {
            data, err := os.ReadFile(path + file.Name())
            if err != nil {
                return nil, err
            }

            var schema Schema
            if err := json.Unmarshal(data, &schema); err != nil {
                return nil, err
            }

            schemas = append(schemas, schema)
            fmt.Printf("✅ Loaded schema: %s\n", schema.Model)
        }
    }

    return schemas, nil
}
