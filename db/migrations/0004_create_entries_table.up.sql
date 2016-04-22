CREATE TABLE entries(
	id 			INT PRIMARY KEY,
	series_id 	INT REFERENCES series(id),
	position    INT REFERENCES digits(position),
	status      BOOLEAN,
	timestamp 	TIMESTAMP DEFAULT now()
)