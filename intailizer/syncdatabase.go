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

	//alter the username
	// alter := `alter table cometoseeauth
	// add column username text not null`

	// _, err = DB.Exec(alter)

	// if err != nil {
	// 	log.Printf("errror altering table")
	// }

	messgingtable := `
	create table if not exists messagetable(
	id serial primary key,
	sender text not null,
	room text not null,
	message text not null,
	sent_at timestamp not null default now()
	)
	`
	_, err = DB.Exec(messgingtable)

	if err != nil {
		log.Fatalf("error in creting messaging table :%v", err)
	}

	fmt.Print("table ready")
}
