package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserItem struct {
	Id        int64  `json:"id"`
	Uid       string `json:"uid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AdminUsersResp struct {
	List  []AdminUserItem `json:"list"`
	Total int64           `json:"total"`
}

func listAdminUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize

		where := "1=1"
		args := make([]any, 0)
		if keyword != "" {
			where = "(username LIKE ? OR email LIKE ?)"
			kw := "%" + keyword + "%"
			args = append(args, kw, kw)
		}

		var total int64
		countSQL := fmt.Sprintf("SELECT COUNT(1) FROM user WHERE %s", where)
		if err := db.QueryRow(countSQL, args...).Scan(&total); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		querySQL := fmt.Sprintf(`
SELECT id, username, COALESCE(email,''), COALESCE(role,'user'), COALESCE(status,'active'), created_at, updated_at
FROM user
WHERE %s
ORDER BY id DESC
LIMIT ? OFFSET ?`, where)
		args = append(args, pageSize, offset)

		rows, err := db.Query(querySQL, args...)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer rows.Close()

		list := make([]AdminUserItem, 0)
		for rows.Next() {
			var item AdminUserItem
			var createdAt, updatedAt []byte
			if err := rows.Scan(&item.Id, &item.Username, &item.Email, &item.Role, &item.Status, &createdAt, &updatedAt); err != nil {
				continue
			}
			item.Uid = fmt.Sprintf("UID%03d", item.Id)
			item.CreatedAt = string(createdAt)
			item.UpdatedAt = string(updatedAt)
			// Keep compatibility with simplified admin logic.
			if item.Id == 1 {
				item.Role = "admin"
			}
			list = append(list, item)
		}

		httpx.OkJsonCtx(r.Context(), w, AdminUsersResp{
			List:  list,
			Total: total,
		})
	}
}

type AdminCreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordEncrypted string `json:"passwordEncrypted"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

func createAdminUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AdminCreateUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid body")
			return
		}
		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)
		req.Role = strings.TrimSpace(req.Role)
		req.Status = strings.TrimSpace(req.Status)
		if strings.TrimSpace(req.PasswordEncrypted) != "" {
			p, err := decryptPasswordEncrypted(strings.TrimSpace(req.PasswordEncrypted))
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid password")
				return
			}
			req.Password = p
		}

		if req.Username == "" || req.Password == "" {
			writeJSONError(w, http.StatusBadRequest, "username and password are required")
			return
		}
		if req.Role == "" {
			req.Role = "user"
		}
		if req.Status == "" {
			req.Status = "active"
		}
		if !isValidUserRole(req.Role) {
			writeJSONError(w, http.StatusBadRequest, "invalid role")
			return
		}
		if !isValidUserStatus(req.Status) {
			writeJSONError(w, http.StatusBadRequest, "invalid status")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		res, err := db.Exec("INSERT INTO user (username, password, email, role, status, balance, character_count) VALUES (?, ?, ?, ?, ?, 0, 0)",
			req.Username,
			string(hashedPassword),
			req.Email,
			req.Role,
			req.Status,
		)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id, _ := res.LastInsertId()

		actor := mustUserID(r.Context())
		writeAdminLog(db, actor, fmt.Sprintf("创建了用户 UID%03d（%s）", id, req.Username), getClientIP(r))

		httpx.OkJsonCtx(r.Context(), w, map[string]any{
			"id": id,
		})
	}
}

type AdminUpdateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	Password string `json:"password"`
	PasswordEncrypted string `json:"passwordEncrypted"`
}

func updateAdminUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseIDFromPath("/api/admin/users/", r.URL.Path)
		if !ok || id <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		if id == 1 {
			writeJSONError(w, http.StatusBadRequest, "cannot modify the default admin user")
			return
		}

		var req AdminUpdateUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid body")
			return
		}
		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)
		req.Role = strings.TrimSpace(req.Role)
		req.Status = strings.TrimSpace(req.Status)
		if strings.TrimSpace(req.PasswordEncrypted) != "" {
			p, err := decryptPasswordEncrypted(strings.TrimSpace(req.PasswordEncrypted))
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid password")
				return
			}
			req.Password = p
		}

		sets := make([]string, 0)
		args := make([]any, 0)

		if req.Username != "" {
			sets = append(sets, "username = ?")
			args = append(args, req.Username)
		}
		if req.Email != "" {
			sets = append(sets, "email = ?")
			args = append(args, req.Email)
		}
		if req.Role != "" {
			if !isValidUserRole(req.Role) {
				writeJSONError(w, http.StatusBadRequest, "invalid role")
				return
			}
			sets = append(sets, "role = ?")
			args = append(args, req.Role)
		}
		if req.Status != "" {
			if !isValidUserStatus(req.Status) {
				writeJSONError(w, http.StatusBadRequest, "invalid status")
				return
			}
			sets = append(sets, "status = ?")
			args = append(args, req.Status)
		}
		if strings.TrimSpace(req.Password) != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			sets = append(sets, "password = ?")
			args = append(args, string(hashedPassword))
		}

		if len(sets) == 0 {
			writeJSONError(w, http.StatusBadRequest, "no updates")
			return
		}

		sqlStr := fmt.Sprintf("UPDATE user SET %s WHERE id = ?", strings.Join(sets, ", "))
		args = append(args, id)
		res, err := db.Exec(sqlStr, args...)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			writeJSONError(w, http.StatusNotFound, "not found")
			return
		}

		actor := mustUserID(r.Context())
		writeAdminLog(db, actor, fmt.Sprintf("更新了用户 UID%03d 的信息", id), getClientIP(r))

		httpx.OkJsonCtx(r.Context(), w, map[string]any{"ok": true})
	}
}

func deleteAdminUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseIDFromPath("/api/admin/users/", r.URL.Path)
		if !ok || id <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		if id == 1 {
			writeJSONError(w, http.StatusBadRequest, "cannot delete the default admin user")
			return
		}

		res, err := db.Exec("DELETE FROM user WHERE id = ?", id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			writeJSONError(w, http.StatusNotFound, "not found")
			return
		}

		actor := mustUserID(r.Context())
		writeAdminLog(db, actor, fmt.Sprintf("删除了用户 UID%03d", id), getClientIP(r))

		httpx.OkJsonCtx(r.Context(), w, map[string]any{"ok": true})
	}
}

func isValidUserRole(role string) bool {
	switch role {
	case "admin", "engineer", "user":
		return true
	default:
		return false
	}
}

func isValidUserStatus(status string) bool {
	switch status {
	case "active", "disabled":
		return true
	default:
		return false
	}
}
