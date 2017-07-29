package handlers

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB
func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/queue?sslmode=disable")

	if err != nil {
		log.Print(err)
		panic(err)
	}
}



func insertToken(j *OutMsgJson, tx *sql.Tx) (int, error) {
	var tokenId int
	tokenRequest := `select token_id from token where access_token = $1`
	err := tx.QueryRow(tokenRequest, j.AccessToken).Scan(&tokenId)

	if err == sql.ErrNoRows {
		tokenString := `insert into token(access_token) VALUES($1) RETURNING token_id`
		err = tx.QueryRow(tokenString, j.AccessToken).Scan(&tokenId)
		if err != nil {
			log.Print(err)
			return 0, err
		}

	} else if err != nil {
		return 0, err
	}

	return tokenId, nil

}


func insertEvent(j *OutMsgJson, tx *sql.Tx) (int, error) {
	var eventId int
	eventRequest := `select event_id from event where event_code = $1`
	err := tx.QueryRow(eventRequest, j.EventCode).Scan(&eventId)

	if err == sql.ErrNoRows {
		eventString := `insert into event(event_code) VALUES($1) RETURNING event_id`
		err = tx.QueryRow(eventString, j.EventCode).Scan(&eventId)
		if err != nil {
			log.Print(err)
			return 0, err
		}

	} else if err != nil {
		return 0, err
	}

	return eventId, nil

}

func insertStream(j *OutMsgJson, tx *sql.Tx) (int, error) {
	var streamId int
	streamRequest := `select stream_id from stream where stream_type = $1`
	err := tx.QueryRow(streamRequest, j.StreamType).Scan(&streamId)

	if err == sql.ErrNoRows {
		streamString := `insert into stream(stream_type) VALUES($1) RETURNING stream_id`
		err = tx.QueryRow(streamString, j.StreamType).Scan(&streamId)
		if err != nil {
			log.Print(err)
			return 0, err
		}

	} else if err != nil {
		return 0, err
	}

	return streamId, nil

}

func InsertMessage(j *OutMsgJson, b []byte) error {

	tx, err := db.Begin()
	defer tx.Rollback()
	tokenId, err := insertToken(j, tx)

	if err != nil {
		return err
	}
	eventId, err := insertEvent(j, tx)
	if err != nil {
		return err
	}

	streamId, err := insertStream(j, tx)
	if err != nil {
		return err
	}

	var msgId int

	msgString := `insert into msg(token_id, event_id, stream_id, to_, data) VALUES($1, $2, $3, $4, $5) RETURNING msg_id`
	err = tx.QueryRow(msgString, tokenId, eventId, streamId, j.To, b).Scan(&msgId)
	if err != nil {
		log.Print(err)
		return err
	}

	tx.Commit()

	return nil

}