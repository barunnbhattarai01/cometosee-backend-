package intailizer

import (
	"fmt"
	"log"
)

func Syncdatabase() {
	authtable := `
 create table if not exists cometoseeauth(
 email text primary key,
 password text not null
 )  
   `

	_, err := DB.Exec(authtable)

	if err != nil {
		log.Fatalf("errror in creating auth table %v", err)
	}

	fmt.Print("table ready")
}
