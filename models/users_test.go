package models

import "testing"

func testingUserService() (*UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		users    = "bond"
		password = "password"
		dbname   = "lenslocked_dev"
	)
	psqlInfo := "postgres://bond:password@localhost/lenslocked_dev?sslmode=disable"
	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)
	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "butt",
		Email: "head@email.com",
	}
	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Received %d ", user.ID)
	}

}
