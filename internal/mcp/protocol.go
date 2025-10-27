package mcp

type Request struct {
	Tool string         `json:"tool"`
	Args map[string]any `json:"args"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
