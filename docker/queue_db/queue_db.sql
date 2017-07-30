CREATE TABLE IF NOT EXISTS token (
	token_id SERIAL PRIMARY KEY,
	access_token varchar(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS event (
	event_id SERIAL PRIMARY KEY,
	event_code varchar(255) NOT NULL
);


CREATE TABLE IF NOT EXISTS stream (
	stream_id SERIAL PRIMARY KEY,
	stream_type varchar(64) NOT NULL
);



CREATE TABLE IF NOT EXISTS msg (
	msg_id SERIAL PRIMARY KEY,
	token_id SERIAL references token (token_id),
	event_id SERIAL references event (event_id),
	stream_id SERIAL references stream(stream_id),
	to_ varchar(128) NOT NULL,
	data json NOT NULL
);




insert into stream (stream_type) VALUES('email');
insert into stream (stream_type) VALUES('sms');
insert into stream (stream_type) VALUES('push');
