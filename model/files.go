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

// FilesSend - sending struct
type FilesSend struct {
	ID          string `json:"id"`
	Question    string `json:"question"`
	Slug        string `json:"slug"`
	CreatedAt   string `json:"createdAt"`
	Username    string `json:"username"`
	Unique_Name string `json:"uniqueName"`
}

// LikeModel - get like
type LikeModel struct {
	question_id string `json:"questionId"`
	user_id     string `json:"userId"`
	likes       bool   `json:"like"`
	dislike     bool   `json::dislike"`
}

// GetQuestions - hold questions struct
type GetQuestions struct {
	ID         string `json:"id"`
	Question   string `json:"question"`
	Poster     string `json:"poster"`
	Slug       string `json:"slug"`
	Created_At string `json:"createdAt"`
	Answer     string `json:"answer"`
	Likes      int    `json:"likes"`
}

// GetAnswers - hold all the answer struct
type GetAnswers struct {
	Question_ID string `json:"questionId"`
	Answer      string `json:"answer"`
	Created_At  string `json:"createdAt"`
	Username    string `json:"username"`
	Unique_Name string `json:"uniqueName"`
}
