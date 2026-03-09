CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    nom TEXT,
    prenom TEXT,
    age INT,
    email TEXT
);


insert into users (nom, prenom, age, email) values 
    ('Doe', 'John', 30, 'john.doe@example.com'),
    ('Smith', 'Jane', 25, 'jane.smith@example.com'),
    ('Brown', 'Charlie', 35, 'charlie.brown@example.com');