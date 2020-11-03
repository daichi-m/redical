package main

// Auth is an action that changes the password of the redis connection and refreshes it.
func Auth(password string) error {
	if err := global.ModifyConfig(&DBConfig{password: password}); err != nil {
		return err
	}
	return nil
}
