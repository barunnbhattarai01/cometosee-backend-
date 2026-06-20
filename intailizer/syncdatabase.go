package intailizer

import (
	"fmt"
	"log"
)

func Syncdatabase() {

	//extension for psotgis
	_, err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS postgis;`)
	if err != nil {
		log.Fatalf("error enabling postgis: %v", err)
	}

	authtable := `
 create table if not exists cometoseeauth(
 auth_id serial primary key,
 email text unique not null,
 username text unique not null,
 password text not null
 )  
   `

	_, err = DB.Exec(authtable)

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

	_, err = DB.Exec(`
ALTER TABLE location
ADD COLUMN IF NOT EXISTS geom geography(Point, 4326);
`)
	if err != nil {
		log.Fatalf("error adding geom column: %v", err)
	}

	//creating index to make searcg more fast
	_, err = DB.Exec(`
CREATE INDEX IF NOT EXISTS idx_location_geom
ON location
USING GIST (geom);
`)
	if err != nil {
		log.Fatalf("error creating geom index: %v", err)
	}

	_, err = DB.Exec(`
CREATE OR REPLACE FUNCTION update_geom_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326)::geography;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`)
	if err != nil {
		log.Fatalf("error creating trigger function: %v", err)
	}

	_, err = DB.Exec(`
DROP TRIGGER IF EXISTS set_geom ON location;

CREATE TRIGGER set_geom
BEFORE INSERT OR UPDATE ON location
FOR EACH ROW
EXECUTE FUNCTION update_geom_column();
`)
	if err != nil {
		log.Fatalf("error creating trigger: %v", err)
	}

	//post
	posttable := `create table if not exists post(
   post_id serial primary key,
     auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
   caption text not null,
   images_url text ,
   venue text not null,
   longitude double precision not null,
   latitude double precision not null,
   sport text not null,
   created_at TIMESTAMP DEFAULT NOW()
 );`

	_, err = DB.Exec(posttable)
	if err != nil {
		log.Fatalf("error creating post:%v", err)
	}

	commenttable := `create table if not exists comments(
	 comment_id serial primary key,
	 auth_id int not null REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
	 post_id int not null REFERENCES post(post_id) ON DELETE CASCADE,
	 comment text not null,
	 created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err = DB.Exec(commenttable)
	if err != nil {
		log.Fatalf("error creating comment table:%v", err)
	}

	liketable := `create table if not exists post_likes(
	 like_id serial primary key,
	 created_at TIMESTAMP DEFAULT NOW(),
	 auth_id int not null REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
		 post_id int not null REFERENCES post(post_id) ON DELETE CASCADE,
		 unique(post_id,auth_id)
	);`

	_, err = DB.Exec(liketable)
	if err != nil {
		log.Fatal("error in creating like table")
	}
	//sharepost
	sharepost := `CREATE TABLE if not exists post_shares (
    share_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(post_id, auth_id)
);`

	_, err = DB.Exec(sharepost)
	if err != nil {
		log.Fatal("error in creating share table")
	}

	//slot
	slottable := `CREATE TABLE IF NOT EXISTS post_slots (
    slot_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,

    max_participants INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),

    CHECK (end_time > start_time)
);`

	_, err = DB.Exec(slottable)
	if err != nil {
		log.Fatal("error in creating slot table")
	}

	//slot participant
	slotparticipant := `CREATE TABLE IF NOT EXISTS slot_participants (
    id SERIAL PRIMARY KEY,
    slot_id INT NOT NULL REFERENCES post_slots(slot_id) ON DELETE CASCADE,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,

    joined_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(slot_id, auth_id)
);`

	_, err = DB.Exec(slotparticipant)
	if err != nil {
		log.Fatal("error in creating slot participant table")
	}

	fmt.Print("table ready")
}
