package models

type Engineer struct {
	BaseModel
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
}
