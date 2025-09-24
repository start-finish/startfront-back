package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apigen",
	Short: "API Generator for msgId-based Go projects",
}

// ---------------- INIT COMMAND ----------------
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new msgId-based API project (with Air hot-reload)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if err := createProjectStructure(projectName); err != nil {
			fmt.Println("❌ Error initializing project:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Project initialized at", projectName)
		fmt.Println("\nNext steps:")
		fmt.Println("  cd", projectName)
		fmt.Println("  go mod init", projectName)
		fmt.Println("  go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/postgres")
		fmt.Println("  go install github.com/cosmtrek/air@latest")
		fmt.Println("  export PATH=$PATH:$(go env GOPATH)/bin")
		fmt.Println("  air   # hot-reloads server on save")
	},
}

func createProjectStructure(name string) error {
	// Create necessary directories
	dirs := []string{
		filepath.Join(name, "cmd", "server"),
		filepath.Join(name, "config"),
		filepath.Join(name, "handlers"),
		filepath.Join(name, "models"),
		filepath.Join(name, "router"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}

	// Write cmd/server/main.go
	mainFile := fmt.Sprintf(`package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"%s/config"
	"%s/models"
	"%s/router"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// migrate base models
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	router.SetupRoutes(r, db)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
`, name, name, name)
	if err := os.WriteFile(filepath.Join(name, "cmd", "server", "main.go"), []byte(mainFile), 0o644); err != nil {
		return fmt.Errorf("write main.go: %w", err)
	}

	// Write config/config.go
	configFile := `package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123 dbname=testdb port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Model Home => table "home" (not "homes")
		},
	})
}`
	if err := os.WriteFile(filepath.Join(name, "config", "config.go"), []byte(configFile), 0o644); err != nil {
		return fmt.Errorf("write config.go: %w", err)
	}

	// Write models/user.go
	userModel := "package models\n\n" +
		"type User struct {\n" +
		"\tID       uint   `json:\"id\" gorm:\"primaryKey\"`\n" +
		"\tUsername string `json:\"username\"`\n" +
		"\tPassword string `json:\"password\"`\n" +
		"}\n\n" +
		"func (User) TableName() string { return \"user\" }\n"
	if err := os.WriteFile(filepath.Join(name, "models", "user.go"), []byte(userModel), 0o644); err != nil {
		return fmt.Errorf("write models/user.go: %w", err)
	}

	// Write router/router.go
	routerFile := `package router

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Request struct {
	MsgID string          ` + "`json:\"msgId\"`" + `
	Data  json.RawMessage ` + "`json:\"data\"`" + `
}

type MsgHandler func(*gorm.DB, json.RawMessage, *gin.Context)

// Keep empty map on init; create-api will inject entries and add handlers import
var MsgHandlers = map[string]MsgHandler{}

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/api", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "Invalid JSON",
			})
			return
		}
		handler, exists := MsgHandlers[req.MsgID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "Unknown msgId",
			})
			return
		}
		handler(db, req.Data, c)
	})
}
`
	if err := os.WriteFile(filepath.Join(name, "router", "router.go"), []byte(routerFile), 0o644); err != nil {
		return fmt.Errorf("write router.go: %w", err)
	}

	// Write .air.toml
	airToml := `# Air live-reload config
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server"
bin = "./tmp/main"
delay = 1000
exclude_dir = ["tmp", "vendor", ".git"]

[log]
color = true
time = true

[screen]
clear_on_rebuild = true
`
	if err := os.WriteFile(filepath.Join(name, ".air.toml"), []byte(airToml), 0o644); err != nil {
		return fmt.Errorf("write .air.toml: %w", err)
	}

	// Write .gitignore
	gitignore := "tmp/\n*.log\n"
	if err := os.WriteFile(filepath.Join(name, ".gitignore"), []byte(gitignore), 0o644); err != nil {
		return fmt.Errorf("write .gitignore: %w", err)
	}

	return nil
}

// ---------------- CREATE-API COMMAND ----------------
var createApiCmd = &cobra.Command{
	Use:   "create-api [name]",
	Short: "Create new API handler, model, and register in MsgHandlers",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiName := strings.ToLower(args[0])
		if err := createHandler(apiName); err != nil {
			fmt.Println("❌ create handler error:", err)
			os.Exit(1)
		}
		if err := createModel(apiName); err != nil {
			fmt.Println("❌ create model error:", err)
			os.Exit(1)
		}
		if err := updateRouter(apiName); err != nil {
			fmt.Println("❌ update router error:", err)
			os.Exit(1)
		}
		fmt.Println("✅ API", apiName, "created successfully")
	},
}

// -------- generate HANDLERS for CRUD operations in one file ----------
func createHandler(name string) error {
	title := strings.Title(name)

	// module name for imports
	moduleName := readModuleName(".")
	if moduleName == "" {
		if wd, _ := os.Getwd(); wd != "" {
			moduleName = filepath.Base(wd)
		}
	}

	// Create a single file for all CRUD handlers
	handlerCode := fmt.Sprintf(`package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"%s/models"
)

func Handle%[2]sGet(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID     uint   `+"`json:\"id\"`"+`
		Name   string `+"`json:\"name\"`"+`
		Status string `+"`json:\"status\"`"+`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	// ID is required
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Build query dynamically based on available parameters
	var m models.%[2]s
	query := db.Model(&models.%[2]s{}).Where("id = ?", req.ID)

	// Apply dynamic filters (Name, Status, etc.)
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%%"+req.Name+"%%")
	}

	// Fetch the record based on ID and optional filters
	if err := query.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return the record
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "%[1]s API (get by ID or other params)",
		"data":    m,
	})

}

func Handle%[2]sList(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		Page   *int    `+"`json:\"page\"`"+`
		Limit  *int    `+"`json:\"limit\"`"+`
		Search *string `+"`json:\"search\"`"+`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
			return
		}
	}

	page := 1
	limit := 10
	if req.Page != nil && *req.Page > 0 { page = *req.Page }
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 200 { limit = *req.Limit }
	offset := (page-1)*limit

	q := db.Model(&models.%[2]s{})
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		q = q.Where("name ILIKE ?", "%%"+strings.TrimSpace(*req.Search)+"%%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	var items []models.%[2]s
	if err := q.Order("id ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	response := gin.H{
		"code":   "0",
		"status": "Success",
		"message": "%[1]s API (list)",
		"meta": gin.H{"page": page, "limit": limit, "total": total},
	}
	if len(items) > 0 {
		response["data"] = items
	}

	c.JSON(http.StatusOK, response)
}

func Handle%[2]sInsert(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req models.%[2]s
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// Use reflection to check for duplicates dynamically
	requiredFields := []struct {
		FieldName  string
		FieldValue interface{}
	}{
		{"name", req.Name},
		// Add more required fields here, e.g., {"FieldName", req.FieldValue}
	}

	// Check for duplicates on each required field
	for _, field := range requiredFields {
		if field.FieldValue == "" {
			// Skip empty fields
			continue
		}

		var existingRecord models.%[2]s
		if err := db.Where(fmt.Sprintf("%%s = ?", field.FieldName), field.FieldValue).First(&existingRecord).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "1",
				"status":  "error",
				"error":   fmt.Sprintf("duplicate field: %%s", field.FieldName),
			})
			return
		}
	}

	if err := db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   "0",
		"status": "Success",
		"message": "%[1]s API (inserted)",
		"data":   req,
	})
}

func Handle%[2]sUpdate(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID     uint   `+"`json:\"id\"`"+`
		Name   string `+"`json:\"name\"`"+`
	}

	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}

	// ID is required
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	// Fetch the record by ID
	var m models.%[2]s
	if err := db.First(&m, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		m.Name = req.Name
	}

	// Save the updated record
	if err := db.Save(&m).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": err.Error()})
		return
	}

	// Return the updated record
	c.JSON(http.StatusOK, gin.H{
		"code":   "0",
		"status": "Success",
		"message": "%[1]s API (updated)",
		"data":   m,
	})
}

func Handle%[2]sDelete(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	var req struct {
		ID uint `+"`json:\"id\"`"+`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "invalid request data: " + err.Error()})
		return
	}
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "1", "status": "error", "error": "'id' is required"})
		return
	}

	res := db.Delete(&models.%[2]s{}, req.ID)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "1", "status": "error", "error": res.Error.Error()})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": "1", "status": "error", "error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"status":  "Success",
		"message": "startfront API (deleted)",
	})
}
`, moduleName, title, title, title)

	if err := os.WriteFile(filepath.Join("handlers", name+".go"), []byte(handlerCode), 0o644); err != nil {
		return fmt.Errorf("write handlers/%s.go: %w", name, err)
	}
	return nil
}

func createModel(name string) error {
	title := strings.Title(name)
	code := fmt.Sprintf(`package models

// %[1]s is a sample model you can modify or remove.
type %[1]s struct {
	ID   uint   `+"`json:\"id\" gorm:\"primaryKey\"`"+`
	Name string `+"`json:\"name\"`"+`
}

// Force GORM to use the singular table name "%[2]s".
func (%[1]s) TableName() string { return "%[2]s" }
`, title, name)

	if err := os.WriteFile(filepath.Join("models", name+".go"), []byte(code), 0o644); err != nil {
		return fmt.Errorf("write models/%s.go: %w", name, err)
	}
	return nil
}

func updateRouter(name string) error {
	// Define the path for the router.go file
	routerFilePath := filepath.Join("router", "router.go")

	// Read the existing content of router.go
	routerFile, err := os.ReadFile(routerFilePath)
	if err != nil {
		return fmt.Errorf("read router.go: %w", err)
	}
	content := string(routerFile)

	// Ensure handlers import exists using the module name from go.mod
	moduleName := readModuleName(".")
	if moduleName == "" {
		if wd, _ := os.Getwd(); wd != "" {
			moduleName = filepath.Base(wd)
		}
	}
	handlersImport := fmt.Sprintf(`"%s/handlers"`, moduleName)
	if !strings.Contains(content, handlersImport) {
		content = injectImport(content, handlersImport)
	}

	// Capitalize the name for handler registration
	handlerName := strings.Title(name)
	upName := strings.ToUpper(name)

	// Register handlers for each action (get, list, create, update, delete)
	handlerRegistration := fmt.Sprintf(`
	// %s API
	"%s_get" : handlers.Handle%sGet,
	"%s_list" : handlers.Handle%sList,
	"%s_create" : handlers.Handle%sInsert,
	"%s_update" : handlers.Handle%sUpdate,
	"%s_delete" : handlers.Handle%sDelete,
`, upName, upName, handlerName, upName, handlerName, upName, handlerName, upName, handlerName, upName, handlerName)

	// Inject the handler registration code into the MsgHandlers map
	if strings.Contains(content, "var MsgHandlers = map[string]MsgHandler{}") {
		content = strings.Replace(content,
			"var MsgHandlers = map[string]MsgHandler{}",
			"var MsgHandlers = map[string]MsgHandler{\n\t"+handlerRegistration+"\n}",
			1)
	} else {
		const anchor = "var MsgHandlers = map[string]MsgHandler{"
		idx := strings.Index(content, anchor)
		if idx == -1 {
			return fmt.Errorf("MsgHandlers map not found in router.go")
		}
		after := content[idx+len(anchor):]
		closeIdx := strings.Index(after, "}")
		if closeIdx == -1 {
			return fmt.Errorf("could not find end of MsgHandlers map")
		}
		before := content[:idx+len(anchor)]
		middle := "\n\t" + handlerRegistration + after[:closeIdx]
		rest := after[closeIdx:]
		content = before + middle + rest
	}

	// Write the updated content back to the router.go file
	if err := os.WriteFile(routerFilePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write router.go: %w", err)
	}
	return nil
}

func injectImport(content, newImport string) string {
	if !strings.Contains(content, "import (") {
		return strings.Replace(content, "package router",
			"package router\n\nimport (\n\t"+newImport+"\n)\n", 1)
	}
	lines := strings.Split(content, "\n")
	var out []string
	inImport := false
	already := false

	for _, ln := range lines {
		if strings.HasPrefix(ln, "import (") {
			inImport = true
			out = append(out, ln)
			continue
		}
		if inImport {
			if strings.TrimSpace(ln) == ")" {
				if !already {
					out = append(out, "\t"+newImport)
				}
				out = append(out, ln)
				inImport = false
				continue
			}
			if strings.TrimSpace(ln) == newImport {
				already = true
			}
			out = append(out, ln)
			continue
		}
		out = append(out, ln)
	}
	return strings.Join(out, "\n")
}

func readModuleName(dir string) string {
	f, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createApiCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
