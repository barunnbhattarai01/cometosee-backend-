package intailizer

import (
	"fmt"
	"log"
)

func Syncdatabase() {

	authtable := `
 create table if not exists cometoseeauth(
 auth_id serial primary key,
 email text unique not null,
 username text not null,
 password text not null
 )  
   `

	_, err := DB.Exec(authtable)

	if err != nil {
		log.Fatalf("errror in creating auth table %v", err)
	}

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

	videocalltable := `CREATE TABLE if not exists video_call_sessions (
    id BIGSERIAL PRIMARY KEY,
    connection_id BIGINT NOT NULL,
    initiated_by_user_id BIGINT NOT NULL,
    agora_channel_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'initiated',
    started_at TIMESTAMP NULL,
    ended_at TIMESTAMP NULL,
    duration_seconds INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);`

	_, err = DB.Exec(videocalltable)

	if err != nil {
		log.Fatalf("error in creating video call session table :%v", err)
	}

	//userdetail
	userdetailtable := `CREATE TABLE if not exists userdetailinfo (
	user_detail_id SERIAL PRIMARY KEY,
	auth_id INT NOT NULL,
calling_name TEXT NOT NULL,
sport TEXT NOT NULL,
skill TEXT NOT NULL,
avatar TEXT,
bio TEXT,
created_at TIMESTAMP DEFAULT NOW(),
FOREIGN KEY (auth_id) REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE
);`

	_, err = DB.Exec(userdetailtable)

	if err != nil {
		log.Fatalf("error in creating user detail info table :%v", err)
	}

	//location
	locationtable := `CREATE TABLE if not exists location (
	id SERIAL PRIMARY KEY,
	user_detail_id INT NOT NULL,
country TEXT NOT NULL,
city TEXT NOT NULL,
latitude DOUBLE PRECISION NOT NULL,
longitude DOUBLE PRECISION NOT NULL,
FOREIGN KEY (user_detail_id) REFERENCES userdetailinfo(user_detail_id) ON DELETE CASCADE
);`

	_, err = DB.Exec(locationtable)

	if err != nil {
		log.Fatalf("error in creating location table :%v", err)
	}

	fmt.Print("table ready")
}
