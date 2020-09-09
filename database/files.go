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

// GetQuestions - get all questions
func (f *FilesDatabase) GetQuestions() ([]model.FilesQuestion, error) {
	fqs := make([]model.FilesQuestion, 0)

	rows, err := f.conn.Query(context.Background(), "SELECT * FROM question")
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
		fq := model.FilesQuestion{}
		err := rows.Scan(&fq.ID, &fq.Question, &fq.Poster, &fq.Slug, &fq.Created_At, &fq.Updated_At)

		if err != nil {
			err = errors.New("an error occured")
			return fqs, err
		}

		fqs = append(fqs, fq)
	}

	return fqs, nil
}

// AddPost - add posts to the database
func (f *FilesDatabase) PostQuestion(p model.FilesQuestion) error {
	_, err := f.conn.Exec(context.Background(), "INSERT INTO question (id, question, poster, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", p.ID, p.Question, p.Poster, p.Slug, p.Created_At, p.Updated_At)

	return err
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

// EditAnswer - edits the quesiton
func (f *FilesDatabase) EditAnswer(s string, na string) error {
	t := time.Now().UTC().Format(time.RFC3339)
	_, err := f.conn.Exec(context.Background(), "UPDATE answer SET answer=$1, updated_at=$2 WHERE question_id=$3", na, t, s)

	return err
}

// GetAnswers - get all answers of question
func (f *FilesDatabase) GetAnswers(s string) ([]model.FilesComment, error) {
	fcs := make([]model.FilesComment, 0)

	rows, err := f.conn.Query(context.Background(), "SELECT * FROM answer WHERE question_id=$1", s)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("no answer found")
		return fcs, err
	case err != nil:
		err = errors.New("try again")
		return fcs, err
	}

	for rows.Next() {
		fc := model.FilesComment{}
		err := rows.Scan(&fc.Question_ID, &fc.Answer, &fc.Commenter, &fc.Created_At, &fc.Updated_At)

		if err != nil {
			err = errors.New("an error occured")
			return fcs, err
		}

		fcs = append(fcs, fc)

	}

	return fcs, nil

}
