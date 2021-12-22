package main

type specification struct {
	OpenAPI        string  `json:"openapi,omitempty"`
	Title          string  `json:"title,omitempty"`
	Server         []string `json:"server"`
	Description    string  `json:"description,omitempty"`
	TermsOfService string  `json:"terms_of_service,omitempty"`
	Contact        contact `json:"contact"`
	Paths          map[string]map[string]path `json:"paths"`
}

type contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type path struct {
	Tags []string
}
