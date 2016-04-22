CREATE TABLE series(
	id 			INT PRIMARY KEY,
	user_id 	INT REFERENCES users(id),
	name 		TEXT,
	description TEXT,
	created 	TIMESTAMP DEFAULT now()
)