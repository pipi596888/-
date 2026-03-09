package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"tts-backend/user-api/internal/config"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/golang-jwt/jwt/v4"
)

type Work struct {
	TaskId    string `json:"taskId"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Format    string `json:"format"`
	CreatedAt string `json:"createdAt"`
}

type WorksResp struct {
	List  []Work `json:"list"`
	Total int64  `json:"total"`
}

type User struct {
	Id              int64   `json:"id"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	Balance         float64 `json:"balance"`
	CharacterCount  int64   `json:"characterCount"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResp struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func getWorksHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("mysql", c.Mysql.DataSource)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT task_id, status, progress, format, created_at FROM tts_task ORDER BY created_at DESC LIMIT 20")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer rows.Close()

		var works []Work
		for rows.Next() {
			var work Work
			var createdAt []byte
			if err := rows.Scan(&work.TaskId, &work.Status, &work.Progress, &work.Format, &createdAt); err != nil {
				continue
			}
			work.CreatedAt = string(createdAt)
			works = append(works, work)
		}

		httpx.OkJsonCtx(r.Context(), w, WorksResp{
			List:  works,
			Total: int64(len(works)),
		})
	}
}

func loginHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		db, err := sql.Open("mysql", c.Mysql.DataSource)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer db.Close()

		var user User
		var passwordHash string
		err = db.QueryRow("SELECT id, username, COALESCE(email,''), COALESCE(balance,0), COALESCE(character_count,0), password FROM user WHERE username = ?", req.Username).Scan(
			&user.Id, &user.Username, &user.Email, &user.Balance, &user.CharacterCount, &passwordHash,
		)

		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 401, "message": "用户名或密码错误"})
			return
		}
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if req.Password != passwordHash {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 401, "message": "用户名或密码错误"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": user.Id,
			"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(c.JwtSecret))

		httpx.OkJsonCtx(r.Context(), w, LoginResp{
			Token: tokenString,
			User:  user,
		})
	}
}

func registerHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		db, err := sql.Open("mysql", c.Mysql.DataSource)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer db.Close()

		// 检查用户名是否已存在
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", req.Username).Scan(&exists)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if exists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"code": 400, "message": "用户名已存在"})
			return
		}

		// 加密密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 插入用户
		result, err := db.Exec("INSERT INTO user (username, password, email, balance, character_count) VALUES (?, ?, ?, 0, 0)",
			req.Username, string(hashedPassword), req.Email)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		userId, _ := result.LastInsertId()

		// 生成 token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId":   userId,
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(c.JwtSecret))

		httpx.OkJsonCtx(r.Context(), w, RegisterResp{
			Token: tokenString,
			User: User{
				Id:              userId,
				Username:        req.Username,
				Email:           req.Email,
				Balance:         0,
				CharacterCount:  0,
			},
		})
	}
}
