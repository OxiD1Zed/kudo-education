package domain

type Note struct {
	Id     uint   `json:"id"`
	IdUser uint   `json:"id_user"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
