package main

// Auth is an action that changes the password of the redis connection and refreshes it.
func Auth(password string) error {
	if err := global.ModifyConfig(&DBConfig{password: password}); err != nil {
		return err
	}
	return nil
}

// Select action changes the database of the redis connection and refreshes it.
func Select(db int) error {
	if err := global.ModifyConfig(&DBConfig{database: db}); err != nil {
		return err
	}
	return nil
}
