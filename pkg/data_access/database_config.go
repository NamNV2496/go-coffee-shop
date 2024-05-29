package database

type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     int32
	Database string
}

// func NewDatabaseConfig() DatabaseConfig {

// 	return DatabaseConfig{
// 		Username: "root",
// 		Password: "root",
// 		Host:     "localhost",
// 		Port:     3306,
// 		Database: "coffee",
// 	}
// }
