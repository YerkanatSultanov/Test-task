CREATE TABLE IF NOT EXISTS persons
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) not null,
    surname     VARCHAR(255) not null,
    patronymic  VARCHAR(255),
    age         INTEGER,
    gender      VARCHAR(255),
    nationality VARCHAR(255)
);