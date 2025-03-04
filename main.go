package main

import (
	"fmt"
	"github.com/shaneplunkett/gator/internal/config"
)

func main() {
	var user = "updated name"
	f, _ := config.Read()
	if err := f.SetUser(user); err != nil {
		return
	}
	r, _ := config.Read()
	fmt.Printf("User: %v", r.CurrentUserName)

}
