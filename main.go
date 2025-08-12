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
	// Create folders
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

	// cmd/server/main.go
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

	// config/config.go
	configFile := `package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123 dbname=testdb port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}`
	if err := os.WriteFile(filepath.Join(name, "config", "config.go"), []byte(configFile), 0o644); err != nil {
		return fmt.Errorf("write config.go: %w", err)
	}

	// models/user.go
	userModel := "package models\n\n" +
		"type User struct {\n" +
		"\tID       uint   `json:\"id\" gorm:\"primaryKey\"`\n" +
		"\tUsername string `json:\"username\"`\n" +
		"\tPassword string `json:\"password\"`\n" +
		"}\n"
	if err := os.WriteFile(filepath.Join(name, "models", "user.go"), []byte(userModel), 0o644); err != nil {
		return fmt.Errorf("write models/user.go: %w", err)
	}

	// router/router.go
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

	// .air.toml
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

	// .gitignore
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

func createHandler(name string) error {
	title := strings.Title(name)
	upper := strings.ToUpper(name)

	// Handler template: returns status error on invalid JSON, statusOK otherwise
	code := fmt.Sprintf(`package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Handle%[1]s(db *gorm.DB, data json.RawMessage, c *gin.Context) {
	// Example: validate JSON body (optional)
	var payload map[string]interface{}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid request data: " + err.Error(),
			})
			return
		}
	}

	// TODO: add your DB logic with 'db' here

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "%[2]s API working",
	})
}
`, title, upper)

	if err := os.WriteFile(filepath.Join("handlers", name+".go"), []byte(code), 0o644); err != nil {
		return fmt.Errorf("write handlers/%s.go: %w", name, err)
	}
	return nil
}

func createModel(name string) error {
	title := strings.Title(name)
	code := "package models\n\n" +
		"// " + title + " is a sample model you can modify or remove.\n" +
		"type " + title + " struct {\n" +
		"\tID   uint   `json:\"id\" gorm:\"primaryKey\"`\n" +
		"\tName string `json:\"name\"`\n" +
		"}\n"
	if err := os.WriteFile(filepath.Join("models", name+".go"), []byte(code), 0o644); err != nil {
		return fmt.Errorf("write models/%s.go: %w", name, err)
	}
	return nil
}

func updateRouter(name string) error {
	path := filepath.Join("router", "router.go")
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(b)

	// ensure handlers import exists using module name from go.mod
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

	// inject into MsgHandlers map
	mapLine := fmt.Sprintf(`"%s": handlers.Handle%s,`, strings.ToUpper(name), strings.Title(name))

	if strings.Contains(content, "var MsgHandlers = map[string]MsgHandler{}") {
		content = strings.Replace(content,
			"var MsgHandlers = map[string]MsgHandler{}",
			"var MsgHandlers = map[string]MsgHandler{\n\t"+mapLine+"\n}",
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
		middle := "\n\t" + mapLine + after[:closeIdx]
		rest := after[closeIdx:]
		content = before + middle + rest
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
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

// ---------------- MAIN ENTRY ----------------
func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createApiCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
