package main

import (
	"database/sql"
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shaneplunkett/gator/internal/config"
	"github.com/shaneplunkett/gator/internal/database"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dburl := os.Getenv("DATABASE_URL")
	log.Infof("Database URL: %v", dburl)
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	dbQueries := database.New(db)

	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error Reading Config: %v", err)
	}
	s := &state{db: dbQueries, config: config}
	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]")
	}

	cmdarg := args[0]
	argList := args[1:]
	comm := command{name: cmdarg, arguements: argList}
	err = cmds.run(s, comm)
	if err != nil {
		log.Fatal(err)
	}
}
