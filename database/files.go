package database

import (
	"context"
	"errors"
	"time"

	"github.com/Hamaiz/go-rest-eg/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// FilesDatabase - hold connection
type FilesDatabase struct {
	conn *pgxpool.Pool
}

// NewFilesDatabase - new database
func NewFilesDatabase(conn *pgxpool.Pool) *FilesDatabase {
	return &FilesDatabase{conn}
}

// GetSearchedQuestions - gets all the searched question
func (f *FilesDatabase) GetSearchedQuestions(l string) ([]model.GetQuestions, error) {
	fqs := make([]model.GetQuestions, 0)

	//rows, err := f.conn.Query(context.Background(), `SELECT * FROM question WHERE question similar to '%(' || $1 || ')%'`, l)

	rows, err := f.conn.Query(context.Background(), "SELECT question.id, question.question, question.poster, question.Slug, question.created_at, answer.answer FROM question LEFT JOIN answer ON question.id=answer.question_id WHERE question.question similar to '%(' || $1 || ')%'", l)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no questions found")
		return fqs, err
	case err != nil:
		err = errors.New("an error occured")
		return fqs, err
	}

	defer rows.Close()

	for rows.Next() {
		fq := model.GetQuestions{}
		var answer interface{}

		err := rows.Scan(&fq.ID, &fq.Question, &fq.Poster, &fq.Slug, &fq.Created_At, &answer)

		if err != nil {

			err = errors.New("an error occured")
			return fqs, err
		}

		var likes int
		likes, err = f.GetLikes(fq.ID)
		if err != nil {
			return fqs, err
		}

		if answer == nil {
			fq.Answer = "not answered yet"
		} else {
			fq.Answer = answer.(string)
		}

		fq.Likes = likes

		fqs = append(fqs, fq)
	}

	return fqs, nil

}

// GetQuestions - get all questions
func (f *FilesDatabase) GetQuestions() ([]model.GetQuestions, error) {
	fqs := make([]model.GetQuestions, 0)

	rows, err := f.conn.Query(context.Background(), "SELECT question.id, question.question, question.poster, question.Slug, question.created_at, answer.answer FROM question LEFT JOIN answer ON question.id=answer.question_id")
	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no questions found")
		return fqs, err
	case err != nil:
		err = errors.New("an error occured")
		return fqs, err
	}

	defer rows.Close()

	for rows.Next() {
		fq := model.GetQuestions{}
		var answer interface{}
		err := rows.Scan(&fq.ID, &fq.Question, &fq.Poster, &fq.Slug, &fq.Created_At, &answer)

		if err != nil {
			err = errors.New("an error occured")
			return fqs, err
		}

		var likes int
		likes, err = f.GetLikes(fq.ID)
		if err != nil {
			return fqs, err
		}

		if answer == nil {
			fq.Answer = "not answered yet"
		} else {
			fq.Answer = answer.(string)
		}

		fq.Likes = likes
		fqs = append(fqs, fq)
	}

	return fqs, nil
}

// PostQuestion - add posts to the database
func (f *FilesDatabase) PostQuestion(p model.FilesQuestion) error {
	_, err := f.conn.Exec(context.Background(), "INSERT INTO question (id, question, poster, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", p.ID, p.Question, p.Poster, p.Slug, p.Created_At, p.Updated_At)

	return err
}

// GetQuest - get only one question
func (f *FilesDatabase) GetQuest(s string) (model.FilesSend, error) {
	fq := model.FilesSend{}

	row := f.conn.QueryRow(context.Background(), "SELECT question.id, question, slug, created_at, username, unique_name FROM question JOIN account ON question.poster=account.id WHERE question.slug=$1", s)
	err := row.Scan(&fq.ID, &fq.Question, &fq.Slug, &fq.CreatedAt, &fq.Username, &fq.Unique_Name)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no question found")
		return fq, err
	case err != nil:
		err = errors.New("try again")
		return fq, err
	}

	return fq, nil
}

// GetQuestion - gets the question by taking in slug
func (f *FilesDatabase) GetQuestion(s string) (model.FilesQuestion, error) {
	fq := model.FilesQuestion{}

	row := f.conn.QueryRow(context.Background(), "SELECT * FROM question WHERE id=$1", s)
	err := row.Scan(&fq.ID, &fq.Question, &fq.Poster, &fq.Slug, &fq.Created_At, &fq.Updated_At)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no question found")
		return fq, err
	case err != nil:
		err = errors.New("try again")
		return fq, err
	}

	return fq, nil
}

// EditQuestion - edits the quesiton
func (f *FilesDatabase) EditQuestion(s string, nq string, slug string) error {
	t := time.Now().UTC().Format(time.RFC3339)
	_, err := f.conn.Exec(context.Background(), "UPDATE question SET question=$1, updated_at=$2, slug=$3 WHERE id=$4", nq, t, slug, s)

	if err != nil {
		return err
	}

	return nil
}

// AddAnswer - add answer to the question
func (f *FilesDatabase) AddAnswer(a model.FilesComment) error {
	_, err := f.conn.Exec(context.Background(), "INSERT INTO answer (question_id, answer, commenter, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", a.Question_ID, a.Answer, a.Commenter, a.Created_At, a.Updated_At)

	return err
}

// GetAnswer - get answer from database
func (f *FilesDatabase) GetAnswer(s string, c string) (model.FilesComment, error) {
	fc := model.FilesComment{}

	row := f.conn.QueryRow(context.Background(), "SELECT * FROM answer WHERE question_id=$1 AND commenter=$2", s, c)
	err := row.Scan(&fc.Question_ID, &fc.Answer, &fc.Commenter, &fc.Created_At, &fc.Updated_At)

	return fc, err
}

func (f *FilesDatabase) GetOneAnswer(s string) (string, error) {
	//fc := model.FilesComment{}
	var ans string
	row := f.conn.QueryRow(context.Background(), "SELECT answer FROM answer WHERE question_id=$1", s)

	//err := row.Scan(&fc.Question_ID, &fc.Answer, &fc.Commenter, &fc.Created_At, &fc.Updated_At)
	err := row.Scan(&ans)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no answer found")
		return "", err
	case err != nil:
		err = errors.New("error occured while getting answers")
		return "", err
	}

	return ans, err
}

// EditAnswer - edits the quesiton
func (f *FilesDatabase) EditAnswer(s string, na string) error {
	t := time.Now().UTC().Format(time.RFC3339)
	_, err := f.conn.Exec(context.Background(), "UPDATE answer SET answer=$1, updated_at=$2 WHERE question_id=$3", na, t, s)

	return err
}

// GetAnswers - get all answers of question
func (f *FilesDatabase) GetAnswers(s string) ([]model.GetAnswers, error) {
	fcs := make([]model.GetAnswers, 0)

	rows, err := f.conn.Query(context.Background(), "SELECT answer.question_id, answer.answer, answer.created_at, account.username, account.unique_name FROM answer JOIN account ON answer.commenter=account.id WHERE answer.question_id=$1", s)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no answer found")
		return fcs, err
	case err != nil:
		err = errors.New("try again")
		return fcs, err
	}

	for rows.Next() {
		fc := model.GetAnswers{}
		err := rows.Scan(&fc.Question_ID, &fc.Answer, &fc.Created_At, &fc.Username, &fc.Unique_Name)

		if err != nil {
			err = errors.New("an error occured")
			return fcs, err
		}

		fcs = append(fcs, fc)

	}

	return fcs, nil

}

// Like - add/remove like form the question
func (f *FilesDatabase) Like(id string, u string) error {
	ctx := context.Background()

	var ui bool
	err := f.conn.QueryRow(ctx, "SELECT likes FROM vote WHERE user_id=$1 AND question_id=$2", u, id).Scan(&ui)

	switch {
	case err == pgx.ErrNoRows:
		_, err = f.conn.Exec(ctx, "INSERT INTO vote (question_id, user_id, likes, dislike) VALUES ($1, $2, true, false)", id, u)
		return err
	case err != nil:
		return err
	}

	if ui {
		_, err = f.conn.Exec(ctx, "UPDATE vote SET likes=false, dislike=false WHERE user_id=$1 AND question_id=$2", u, id)

		return err
	}

	_, err = f.conn.Exec(ctx, "UPDATE vote SET likes=true, dislike=false WHERE user_id=$1 AND question_id=$2", u, id)

	return err

}

// Dislike - add/remove like form the question
func (f *FilesDatabase) Dislike(id string, u string) error {
	ctx := context.Background()

	var ui bool
	err := f.conn.QueryRow(ctx, "SELECT dislike FROM vote WHERE user_id=$1 AND question_id=$2", u, id).Scan(&ui)

	switch {
	case err == pgx.ErrNoRows:
		_, err = f.conn.Exec(ctx, "INSERT INTO vote (question_id, user_id, likes, dislike) VALUES ($1, $2, false, true)", id, u)

		return err
	case err != nil:
		return err
	}

	if ui {
		_, err = f.conn.Exec(ctx, "UPDATE vote SET likes=false, dislike=false WHERE user_id=$1 AND question_id=$2", u, id)

		return err
	}

	_, err = f.conn.Exec(ctx, "UPDATE vote SET likes=false, dislike=true WHERE user_id=$1 AND question_id=$2", u, id)

	return err
}

// GetLikes - get all likes of a  post
func (f *FilesDatabase) GetLikes(id string) (int, error) {
	lms := make([]bool, 0)

	rows, err := f.conn.Query(context.Background(), "SELECT likes FROM vote WHERE question_id=$1 AND likes=true", id)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no answer found")
		return 0, err
	case err != nil:
		err = errors.New("try again")
		return 0, err
	}

	for rows.Next() {
		var lm bool
		err := rows.Scan(&lm)

		if err != nil {
			err = errors.New("an error occured")
			return 0, err
		}

		lms = append(lms, lm)
	}

	likes := len(lms)
	return likes, nil

}
