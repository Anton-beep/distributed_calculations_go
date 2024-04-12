package db

type User struct {
	ID       int    `db:"id" json:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

func (a *APIDb) GetUserById(id int) (User, error) {
	user := User{}
	err := a.db.QueryRow("SELECT * FROM users WHERE id=$1", id).
		Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *APIDb) GetUserByUsername(username string) (User, error) {
	user := User{}
	err := a.db.QueryRow("SELECT * FROM users WHERE login=$1", username).
		Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *APIDb) AddUser(user User) error {
	_, err := a.db.Exec("INSERT INTO users(login, password) VALUES ($1, $2)", user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIDb) DeleteUser(id int) error {
	_, err := a.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIDb) UpdateUser(user User) error {
	_, err := a.db.Exec("UPDATE users SET login=$1, password=$2 WHERE id=$3", user.Login, user.Password, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *APIDb) GetLastIdUsers() (int, error) {
	var id int
	err := a.db.QueryRow("SELECT id FROM users ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
