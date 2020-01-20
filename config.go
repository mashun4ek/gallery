package main


type PostgresConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Name string `json:"name"`
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:   "localhost",
		Port: 5432,
		User: "mariaker",
		Password: "nopassword",
		Name: "gallery",
	}
}

# main.go
const (
	
)

isProd := false
s
fmt.Println("Starting the server on :8080...")
http.ListenAndServe(":8080", 

# models/users.go
// pepper
const userPwPepper = "unique8!@gallery!"
const hmacSecretKey = "secret-hmac-key"

# models/services
db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)