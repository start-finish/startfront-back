package engine

type Schema struct {
	Model  string  `json:"model"`
	Fields []Field `json:"fields"`
	Routes struct {
		List   bool `json:"list"`
		Get    bool `json:"get"`
		Create bool `json:"create"`
		Update bool `json:"update"`
		Delete bool `json:"delete"`
	} `json:"routes"`
}

type Field struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	PrimaryKey    bool   `json:"primary_key"`
	AutoIncrement bool   `json:"auto_increment"`
	Required      bool   `json:"required"`
	Index         bool   `json:"index"`
	Unique        bool   `json:"unique"`
}
