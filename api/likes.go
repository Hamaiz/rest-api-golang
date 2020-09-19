package api

import (
	"net/http"
	"strconv"

	"github.com/Hamaiz/go-rest-eg/helper"
)

// LikesHandler - likes the post - @POST - /api/like
func (m *Media) LikesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	// if not logged in
	if !m.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	userId, err := m.store.GetUser(r)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	// get form value
	id := r.FormValue("id")
	if id == "" {
		helper.ASM(w, 403, "an error occured")
		return
	}

	// if question exists
	_, err = m.conn.GetQuestion(id)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	err = m.conn.Like(id, userId)
	if err != nil {
		helper.ASM(w, 403, "error occured while adding like")
		return
	}

	helper.ASM(w, 200, "done")
}

// DislikesHandler - dislikes the post - @POST - /api/dislike
func (m *Media) DislikesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	// if not logged in
	if !m.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// get user id
	userId, err := m.store.GetUser(r)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	// get form value
	id := r.FormValue("id")
	if id == "" {
		helper.ASM(w, 403, "an error occured")
		return
	}

	// if question exists
	_, err = m.conn.GetQuestion(id)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	err = m.conn.Dislike(id, userId)
	if err != nil {
		helper.ASM(w, 403, "error occured while adding dislike")
		return
	}

	helper.ASM(w, 200, "done")
}

// GetLikesHandler - get the number of like - @POST - /api/get-likes
func (m *Media) GetLikesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	// form value
	id := r.FormValue("id")
	if id == "" {
		helper.ASM(w, 404, "an error occured")
		return
	}

	// if question exists
	_, err := m.conn.GetQuestion(id)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	// get likes
	var likes int
	likes, err = m.conn.GetLikes(id)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	w.Write([]byte(`{"likes":` + strconv.Itoa(likes) + `}`))
}
