package repository

import (
	"database/sql"
	"errors"
	"log"

	"20dojo-online/pkg/domain/model"
	"20dojo-online/pkg/domain/repository"
	"20dojo-online/pkg/infra/mysql"
)

type UserRepo interface {
	Create(id, authToken, name string, coin, stage int32) error
	SelectByAuthToken(authToken string) (*model.User, error)
	Update(record *model.User, name string, coin, stage int32) error
	Lock(user *model.User) error
}

type userRepo struct {
	mysql.SqlHandler
}

// NewUserRepo repository/user.goで用意したインタフェースの実装
func NewUserRepo(sqlHandler mysql.SqlHandler) repository.UserRepo {
	userRepo := userRepo{sqlHandler}
	return &userRepo
}

// Create ユーザーの登録
func (ur *userRepo) Create(id, authToken, name string, coin, stage int32) error {
	stmt, err := ur.SqlHandler.Conn.Prepare("INSERT INTO user (id, auth_token, name, coin, stage) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, authToken, name, coin, stage)
	return err
}

// SelectByAuthToken auth_tokenを条件にレコードを取得
func (ur *userRepo) SelectByAuthToken(authToken string) (*model.User, error) {
	row := ur.SqlHandler.Conn.QueryRow("SELECT * FROM user WHERE auth_token = ?", authToken)
	return convertToUser(row)
}

// Update ユーザー情報の更新
func (ur *userRepo) Update(record *model.User, name string, coin, stage int32) error {
	stmt, err := ur.SqlHandler.Conn.Prepare("UPDATE user SET name = ?, coin = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, coin, record.ID)
	return err
}

// Lock ユーザーをロック
func (ur *userRepo) Lock(user *model.User) error {
	row := ur.SqlHandler.Conn.QueryRow("SELECT * FROM user WHERE id = ? FOR UPDATE", user.ID)
	if row == nil {
		return errors.New("user not found")
	}
	return nil
}

// convertToUser rowデータをUserデータへ変換
func convertToUser(row *sql.Row) (*model.User, error) {
	user := model.User{}
	if err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.Coin, &user.Stage); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

// convertToUsers rowsデータをUsersデータへ変換する
func convertToUsers(rows *sql.Rows) (model.Users, error) {
	var users model.Users
	for rows.Next() {
		user := model.User{}
		if err := rows.Scan(&user.ID, &user.AuthToken, &user.Name, &user.Coin, &user.Stage); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}