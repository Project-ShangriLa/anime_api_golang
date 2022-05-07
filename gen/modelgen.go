package main

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// generate code
func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{OutPath: "../model"})

	g.UseDB(gormConnect())

	g.GenerateAllTable()

	g.Execute()
}

func gormConnect() *gorm.DB {
	var err error

	dbHost := os.Getenv("ANIME_API_DB_HOST")
	dbUser := os.Getenv("ANIME_API_DB_USER")
	dbPass := os.Getenv("ANIME_API_DB_PASS")

	if len(dbHost) == 0 {
		dbUser = "root"
	}

	if len(dbUser) > 0 {
		dbPass = ":" + dbPass
	}

	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	db, err := gorm.Open(mysql.Open(dbUser + dbPass + "@" + "tcp(" + dbHost + ")/anime_admin_development?parseTime=true"))
	if err != nil {
		panic(err.Error())
	}

	return db
}
