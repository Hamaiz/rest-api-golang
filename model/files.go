package model

// FilesQuestion - define the arch of question
type FilesQuestion struct {
	ID         string `json:"id"`
	Question   string `json:"question"`
	Poster     string `json:"poster"`
	Slug       string `json:"slug"`
	Created_At string `json:"createdAt"`
	Updated_At string `json:"updatedAt"`
}

// FilesComment - define comment of question
type FilesComment struct {
	Question_ID string `json:"questionId"`
	Answer      string `json:"answer"`
	Commenter   string `json:"commenter"`
	Created_At  string `json:"createdAt"`
	Updated_At  string `json:"updatedAt"`
}
