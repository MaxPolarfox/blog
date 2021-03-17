package types

type Blog struct {
	ID       string `json:"id"`
	AuthorId string `json:"author_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}
