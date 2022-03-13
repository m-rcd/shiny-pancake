package sql

const CreateNoteTable = `
CREATE TABLE if not exists notes (
    id INT unsigned NOT NULL AUTO_INCREMENT, 
    name VARCHAR(150) NOT NULL, 
    content VARCHAR(150) NOT NULL, 
	archived BOOLEAN NOT NULL,
	username VARCHAR(150) NOT NULL,
    PRIMARY KEY     (id)  
    );`
