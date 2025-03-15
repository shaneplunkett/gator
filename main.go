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
		log.Fatalf("Error Loading Environment Config: %v", err)
	}

	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error Reading Config: %v", err)
	}

	dburl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	s := &state{db: dbQueries, config: config}

	cmds := commands{make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	comm := command{name: cmdName, arguments: cmdArgs}
	err = cmds.run(s, comm)
	if err != nil {
		log.Fatal(err)
	}
}
