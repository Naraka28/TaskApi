package user

import (
	"database/sql"
	"fmt"
)

type UserRepository struct{
	r *sql.DB
}

func NewRepository(db *sql.DB) *UserRepository{
	return &UserRepository{ r: db }
}

func (repo *UserRepository) GetAll() ([]User, error){
	var users []User

	rows, err := repo.r.Query("SELECT id, name, age, email FROM users")

	if err != nil {
		return nil, fmt.Errorf("Error gathering users, message: %v", err)
	}
	defer rows.Close()

	for rows.Next(){
		var user User
		if err := rows.Scan(&user.Id,&user.Name, &user.Age, &user.Email); err != nil{
			return nil, fmt.Errorf("Getting users failed, message: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Error: %v", err)
	}
	return users, nil
}
func (repo *UserRepository) GetUserById(id int) (User, error){
	var user User

	row := repo.r.QueryRow("SELECT * FROM users WHERE id = ?", id)

	if err := row.Scan(&user.Id,&user.Name, &user.Age, &user.Email); err != nil{
		if err == sql.ErrNoRows{
			return user, fmt.Errorf("No such User")
		}
		return user, err
	}
	return user, nil
}

func (repo *UserRepository) Register(user RegisterUser) (int64, error) {
	var id int64

	result, err := repo.r.Exec("INSERT INTO users(name, age, email, password) VALUES(?,?,?,?)",user.Name, user.Age, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("Failed to register: %v", err)
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error ID retrieval: %v", err)
	}

	return id, nil
}
func (repo *UserRepository) GetByEmail(email string) (User, error){
	var user User
	result := repo.r.QueryRow("SELECT * FROM users WHERE email = ?", email)

	if err := result.Scan(&user.Id, &user.Name, &user.Age, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows{
			return user, fmt.Errorf("No user with email: %s", email)
		}
		return user, fmt.Errorf("GetByEmail: %v", err)
	}

	return user, nil
}