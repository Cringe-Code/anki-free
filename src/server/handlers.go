package server

import (
	"anki"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user anki.UserReq

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		s.logger.Error("error while parse json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while parse json"}`))
		return
	}

	if user.Login == "" || user.Name == "" || user.Password == "" {
		s.logger.Error("cant create user", "error", "empty register fields")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"reason": "empty register fields"}`))
		return
	}
	var exists bool
	q := "select exists(select 1 from users where login=$1)"
	err = s.db.QueryRow(q, user.Login).Scan(&exists)

	if err != nil {
		s.logger.Error("error while checking exists", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while checking exists"}`))
		return
	}

	if exists {
		s.logger.Error("cant create user", "error", "user with same login exists")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"reason": "user with same login exists"}`))
		return
	}

	hash := sha512.Sum512([]byte(user.Login + user.Password))
	hashedPassword := hex.EncodeToString(hash[:])

	q = "insert into users(name, login, hash_password) values ($1, $2, $3)"
	_, err = s.db.Exec(q, user.Name, user.Login, hashedPassword)

	if err != nil {
		s.logger.Error("error while insert into users table", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while insert into users table"}`))
		return
	}

	ans := anki.UserRes{
		Name: user.Name,
	}

	userJson, _ := json.Marshal(ans)

	s.logger.Info("user has been created")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"profile": %s}`, userJson)))
}

func (s *Server) handlerSignIn(w http.ResponseWriter, r *http.Request) {
	var user anki.UserReq

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		s.logger.Error("error while parse json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while parse json"}`))
		return
	}

	if user.Login == "" || user.Password == "" {

	}
}

func (s *Server) handleCreatePack(w http.ResponseWriter, r *http.Request) {

	// (´｡• ω •｡`)
	var user anki.User
	_ = user
	// тут будет проверка, что пользователь авторизован, но я не крудошлепа, чтобы такое писать

	var pack anki.Pack

	err := json.NewDecoder(r.Body).Decode(&pack)

	if err != nil {
		s.logger.Error("error while parse json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while parse json"}`))
		return
	}

	if pack.Name == "" {
		s.logger.Error("cant create pack", "error", "empty pack create fields")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"reason": "empty pack create fields"}`))
		return
	}

	var exists bool

	q := "select exists(selec 1 from packs where name=$1)"
	err = s.db.QueryRow(q, pack.Name).Scan(&exists)

	if err != nil {
		s.logger.Error("error while checking pack exists", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while checking pack exists"}`))
		return
	}

	if exists {
		s.logger.Error("cant create pack", "error", "pack with same name exists")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"reason": "pack with same name exists"}`))
		return
	}

	q = "insert into packs(name, rank) values ($1, $2) returning id"
	var packId string
	err = s.db.QueryRow(q, pack.Name, pack.Rank).Scan(packId)

	if err != nil {
		s.logger.Error("error while insert into packs table", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while insert into packs table"}`))
		return
	}

	q = "insert into user_pack (user_id, pack_id) values ($1, $2)"
	_, err = s.db.Exec(q, user.Id, packId)

	if err != nil {
		s.logger.Error("error while insert into packs table", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while insert into packs table"}`))
		return
	}

	ans := anki.Pack{
		Name: pack.Name,
		Rank: pack.Rank,
	}

	packJson, _ := json.Marshal(ans)

	s.logger.Info("pack has been created")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"pack": %s}`, packJson)))
}
