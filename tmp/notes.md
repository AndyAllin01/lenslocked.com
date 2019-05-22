
 #   models/users.go

    const userPwPepper = "secret-random-string-this-project"
const hmacSecretKey = "secret-hmac-key"

#models/services.go

	db, err := gorm.Open("postgres", connectionInfo)

	db.LogMode(true)