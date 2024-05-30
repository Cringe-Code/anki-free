package server

import (
	"anki"
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

	hashedPassword := hashPassword(user.Password, user.Login)

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
		s.logger.Error("cant sign in user", "error", "empty sign in fields")
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

	if !exists {
		s.logger.Error("cant sign in user", "error", "user with such login or password doesnt exists")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"reason": "wrong login or password"}`))
		return
	}

	var userHashPassword string
	q = "select hash_password from users where login=$1"

	err = s.db.QueryRow(q, user.Login).Scan(&userHashPassword)

	if err != nil {
		s.logger.Error("error while get user password", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while get user password"}`))
		return
	}

	if userHashPassword != hashPassword(user.Password, user.Login) {
		s.logger.Error("cant sign in user", "error", "user with such login or password doesnt exists")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"reason": "wrong login or password"}`))
		return
	}

	token, err := s.generateToken(user.Login)

	if err != nil {
		s.logger.Error("error while generate token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while generate jwt-token}`))
		return
	}

	TokenJson, _ := json.Marshal(token)
	s.logger.Info("user sign in successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(TokenJson))
}

func (s *Server) handleCreatePack(w http.ResponseWriter, r *http.Request) {

	// (´｡• ω •｡`)

	claims := checkAuth(w, r, s.signingKey)

	if claims == nil {
		return
	}

	var user anki.User

	userLogin, ok := claims["login"].(string)

	if !ok {
		s.logger.Error("cant create pack", "error", "error while parse jwt-claims")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"reason": "error while parse jwt-claims"}`))
		return
	}

	q := "select id, name, login from users where login=$1"

	row := s.db.QueryRow(q, userLogin)

	row.Scan(&user.Id, &user.Name, &user.Login)
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

	q = "select exists(select 1 from packs where name=$1)"
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
	err = s.db.QueryRow(q, pack.Name, pack.Rank).Scan(&packId)

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

func (s *Server) handleAddWord(w http.ResponseWriter, r *http.Request) {

}
