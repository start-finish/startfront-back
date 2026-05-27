package docs

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---- OpenAPI 3.0 Spec Types ----

type OpenAPISpec struct {
	OpenAPI    string                `json:"openapi"`
	Info       Info                  `json:"info"`
	Servers    []Server              `json:"servers,omitempty"`
	Paths      map[string]PathItem   `json:"paths"`
	Components *Components           `json:"components,omitempty"`
	Tags       []Tag                 `json:"tags,omitempty"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type PathItem struct {
	Post *Operation `json:"post,omitempty"`
}

type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

type RequestBody struct {
	Required bool             `json:"required"`
	Content  map[string]Media `json:"content"`
}

type Media struct {
	Schema   SchemaRef          `json:"schema"`
	Examples map[string]Example `json:"examples,omitempty"`
}

type Example struct {
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value"`
}

type Response struct {
	Description string           `json:"description"`
	Content     map[string]Media `json:"content,omitempty"`
}

type SchemaRef struct {
	Ref        string              `json:"$ref,omitempty"`
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Items      *SchemaRef          `json:"items,omitempty"`
	Example    interface{}         `json:"example,omitempty"`
}

type Property struct {
	Type        string      `json:"type,omitempty"`
	Format      string      `json:"format,omitempty"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
	Items       *SchemaRef  `json:"items,omitempty"`
}

type Components struct {
	Schemas map[string]SchemaRef `json:"schemas"`
}

// ---- Spec Builder ----

func buildSpec() OpenAPISpec {
	spec := OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       "Startfront API",
			Description: "All requests go through a single endpoint `POST /api/startProcess` with a JSON body `{\"msgId\": \"...\", \"data\": {...}}`.\n\nUse the **Request Body Examples** dropdown below to select and try different operations (msgIds).",
			Version:     "1.0.0",
		},
		Servers: []Server{
			{URL: "http://localhost:8080", Description: "Local dev server"},
		},
		Paths:      map[string]PathItem{},
		Components: &Components{Schemas: map[string]SchemaRef{}},
	}

	// Initialize the single shared operation
	examples := map[string]Example{}
	spec.Paths["/api/startProcess"] = PathItem{
		Post: &Operation{
			Tags:        []string{"API Processes"},
			Summary:     "Start Process (Master Endpoint)",
			Description: "This single endpoint handles all application processes based on the `msgId` provided in the request body.",
			OperationID: "startProcess",
			RequestBody: &RequestBody{
				Required: true,
				Content: map[string]Media{
					"application/json": {
						Schema: SchemaRef{
							Type: "object",
							Properties: map[string]Property{
								"msgId": {Type: "string", Description: "The unique identifier for the operation"},
								"data":  {Type: "object", Description: "The payload for the operation"},
							},
						},
						Examples: examples,
					},
				},
			},
			Responses: standardResponses(),
		},
	}

	// ─── Resource definitions ───
	type fieldDef = struct {
		name    string
		typ     string
		format  string
		example interface{}
	}

	type resource struct {
		tag    string // display tag (group name)
		prefix string // msgId prefix e.g. "USERS"
		fields []fieldDef
	}

	resources := []resource{
		{tag: "Users", prefix: "USERS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"email", "string", "email", "user@example.com"},
			{"username", "string", "", "johndoe"},
			{"status", "string", "", "active"},
		}},
		{tag: "Auth", prefix: "AUTH"},
		{tag: "Clients", prefix: "CLIENTS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"name", "string", "", "Acme Corp"},
			{"status", "string", "", "active"},
		}},
		{tag: "Projects", prefix: "PROJECTS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"client_id", "integer", "", 1},
			{"name", "string", "", "My Project"},
			{"status", "string", "", "active"},
		}},
		{tag: "Client Users", prefix: "CLIENT_USERS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"client_id", "integer", "", 1},
			{"user_id", "integer", "", 1},
		}},
		{tag: "Screens", prefix: "SCREENS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"name", "string", "", "Home Screen"},
			{"route_path", "string", "", "/home"},
			{"description", "string", "", "Main screen"},
			{"is_active", "string", "", "0"},
			{"created_by", "string", "", "admin"},
		}},
		{tag: "Navigation Menus", prefix: "NAVIGATION_MENUS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"name", "string", "", "Main Menu"},
			{"scope", "string", "", "global"},
			{"client_id", "string", "", ""},
		}},
		{tag: "Navigation Items", prefix: "NAVIGATION_ITEMS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"menu_id", "integer", "", 1},
			{"parent_id", "integer", "", nil},
			{"label", "string", "", "Dashboard"},
			{"path", "string", "", "/dashboard"},
			{"screen_id", "integer", "", nil},
			{"icon", "string", "", "dashboard"},
			{"sort_order", "integer", "", 0},
			{"properties", "object", "", map[string]interface{}{}},
		}},
		{tag: "Widgets", prefix: "WIDGETS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"key", "string", "", "kpi_card"},
			{"label", "string", "", "KPI Card"},
			{"category", "string", "", "display"},
			{"icon_type", "string", "", "flutter"},
			{"icon_value", "string", "", "Icons.dashboard"},
			{"is_builtin", "boolean", "", true},
			{"version", "string", "", "1.0.0"},
			{"config_schema", "object", "", map[string]interface{}{}},
		}},
		{tag: "Widget Instances", prefix: "WIDGET_INSTANCES", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"widget_id", "integer", "", 1},
			{"title", "string", "", "Sales KPI"},
			{"scope", "string", "", "global"},
			{"client_id", "integer", "", nil},
			{"screen_id", "integer", "", nil},
			{"config", "object", "", map[string]interface{}{}},
			{"sort_order", "integer", "", 0},
			{"is_active", "boolean", "", true},
		}},
		{tag: "Screen Widgets", prefix: "SCREEN_WIDGETS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"screen_id", "integer", "", 1},
			{"widget_instance_id", "integer", "", 1},
			{"layout", "object", "", map[string]interface{}{}},
			{"sort_order", "integer", "", 0},
		}},
		{tag: "Widget Presets", prefix: "WIDGET_PRESETS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"name", "string", "", "Default Dashboard"},
			{"description", "string", "", "Default preset"},
			{"scope", "string", "", "global"},
			{"client_id", "integer", "", nil},
			{"created_by", "integer", "", nil},
		}},
		{tag: "Widget Preset Items", prefix: "WIDGET_PRESET_ITEMS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"preset_id", "integer", "", 1},
			{"widget_id", "integer", "", 1},
			{"default_config", "object", "", map[string]interface{}{}},
			{"sort_order", "integer", "", 0},
		}},
		{tag: "Themes", prefix: "THEMES", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"name", "string", "", "Dark Theme"},
			{"scope", "string", "", "global"},
			{"client_id", "integer", "", nil},
			{"variables", "object", "", map[string]interface{}{}},
			{"is_active", "boolean", "", false},
		}},
		{tag: "Settings", prefix: "SETTINGS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"key", "string", "", "site_name"},
			{"scope", "string", "", "global"},
			{"client_id", "integer", "", nil},
			{"value", "object", "", map[string]interface{}{}},
		}},
		{tag: "Analytics Events", prefix: "ANALYTICS_EVENTS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"occurred_at", "string", "date-time", "2026-01-01T00:00:00Z"},
			{"user_id", "integer", "", nil},
			{"client_id", "integer", "", nil},
			{"event_name", "string", "", "page_view"},
			{"screen_id", "integer", "", nil},
			{"widget_instance_id", "integer", "", nil},
			{"meta", "object", "", map[string]interface{}{}},
		}},
		{tag: "Roles", prefix: "ROLES", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"role_name", "string", "", "admin"},
		}},
		{tag: "User Roles", prefix: "USER_ROLES", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"user_id", "integer", "", 1},
			{"role_id", "integer", "", 1},
		}},
		{tag: "Permissions", prefix: "PERMISSIONS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"code", "string", "", "manage_users"},
			{"description", "string", "", "Can manage users"},
		}},
		{tag: "Role Permissions", prefix: "ROLE_PERMISSIONS", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"role_id", "integer", "", 1},
			{"permission_id", "integer", "", 1},
		}},
		{tag: "Admin Dashboard", prefix: "ADMIN_DASHBOARD", fields: []fieldDef{
			{"id", "integer", "", 1},
			{"name", "string", "", "Overview"},
		}},
	}

	// ─── Auth Examples ───
	addExample(examples, "Auth", "login", "Login",
		"Authenticate with username/email and password.",
		map[string]interface{}{
			"msgId": "login",
			"data": map[string]interface{}{
				"email":    "user@example.com",
				"username": "johndoe",
				"password": "secret123",
			},
		},
	)
	addExample(examples, "Auth", "sign_up", "Sign Up",
		"Register a new user.",
		map[string]interface{}{
			"msgId": "sign_up",
			"data": map[string]interface{}{
				"username": "newuser",
				"email":    "new@example.com",
				"password": "password123",
			},
		},
	)

	// ─── Resource Examples ───
	for _, res := range resources {
		if res.prefix == "AUTH" {
			continue
		}

		// add schema to components
		props := map[string]Property{}
		for _, f := range res.fields {
			props[f.name] = Property{Type: f.typ, Format: f.format, Example: f.example}
		}
		spec.Components.Schemas[res.prefix] = SchemaRef{Type: "object", Properties: props}

		// GET (by ID)
		addExample(examples, res.tag, res.prefix+"_get", res.tag+" — Get by ID",
			"Fetch a single record by `id`.",
			map[string]interface{}{"msgId": res.prefix + "_get", "data": map[string]interface{}{"id": 1}},
		)

		// LIST
		addExample(examples, res.tag, res.prefix+"_list", res.tag+" — List all",
			"Paginated list.",
			map[string]interface{}{"msgId": res.prefix + "_list", "data": map[string]interface{}{"page": 1, "limit": 10}},
		)

		// CREATE
		exampleCreate := map[string]interface{}{}
		for _, f := range res.fields {
			if f.name == "id" {
				continue
			}
			if f.example != nil {
				exampleCreate[f.name] = f.example
			}
		}
		addExample(examples, res.tag, res.prefix+"_create", res.tag+" — Create",
			"Insert a new record.",
			map[string]interface{}{"msgId": res.prefix + "_create", "data": exampleCreate},
		)

		// UPDATE
		exampleUpdate := map[string]interface{}{"id": 1}
		for _, f := range res.fields {
			if f.name == "id" {
				continue
			}
			if f.name == "name" || f.name == "label" || f.name == "key" || f.name == "username" {
				exampleUpdate[f.name] = f.example
			}
		}
		addExample(examples, res.tag, res.prefix+"_update", res.tag+" — Update",
			"Update an existing record by `id`.",
			map[string]interface{}{"msgId": res.prefix + "_update", "data": exampleUpdate},
		)

		// DELETE
		addExample(examples, res.tag, res.prefix+"_delete", res.tag+" — Delete",
			"Delete a record by `id`.",
			map[string]interface{}{"msgId": res.prefix + "_delete", "data": map[string]interface{}{"id": 1}},
		)
	}

	return spec
}

func addExample(examples map[string]Example, tag, key, summary, desc string, value interface{}) {
	examples[key] = Example{
		Summary:     "[" + tag + "] " + summary,
		Description: desc,
		Value:       value,
	}
}

func standardResponses() map[string]Response {
	return map[string]Response{
		"200": {
			Description: "Success",
			Content: map[string]Media{
				"application/json": {
					Schema: SchemaRef{
						Type: "object",
						Properties: map[string]Property{
							"code":    {Type: "string", Example: "0"},
							"status":  {Type: "string", Example: "Success"},
							"message": {Type: "string"},
							"data":    {Type: "object"},
						},
					},
				},
			},
		},
		"400": {Description: "Bad request / validation error"},
		"404": {Description: "Record not found"},
		"500": {Description: "Internal server error"},
	}
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Startfront API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *::before, *::after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: '/swagger/doc.json',
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset,
      ],
      layout: 'BaseLayout',
    });
  </script>
</body>
</html>`

func SetupSwagger(r *gin.Engine) {
	spec := buildSpec()

	r.GET("/swagger/doc.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, spec)
	})

	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})

	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	r.GET("/swagger/spec", func(c *gin.Context) {
		data, _ := json.MarshalIndent(spec, "", "  ")
		c.Data(http.StatusOK, "application/json", data)
	})
}
