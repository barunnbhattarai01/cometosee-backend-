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

	subscritiontable := `
create table if not exists subscriptiontable(
id serial primary key,
	user_email text not null,
	start_date timestamp not null default now(),
	end_date timestamp
)
	`

	_, err = DB.Exec(subscritiontable)

	if err != nil {
		log.Fatalf("error in creating subscription table :%v", err)
	}

	paymenttable := `
create table if not exists paymenttable(
id serial primary key,
	user_email text not null,
	amount decimal(10, 2) not null,
	payment_date timestamp not null default now()
)
	`

	_, err = DB.Exec(paymenttable)

	if err != nil {
		log.Fatalf("error in creating payment table :%v", err)
	}

	connectiontable := `CREATE TABLE if not exists connectionstable(
    id BIGSERIAL PRIMARY KEY,
    user_id_1 text NOT NULL,
    user_id_2 text NOT NULL,
    status TEXT NOT NULL,
    requested_by text NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),

    CHECK (user_id_1 < user_id_2)
);

CREATE UNIQUE INDEX if not exists unique_pair
ON connectionstable(user_id_1, user_id_2);`

	_, err = DB.Exec(connectiontable)

	if err != nil {
		log.Fatalf("error in creating connection table :%v", err)
	}

	fmt.Print("table ready")
}
