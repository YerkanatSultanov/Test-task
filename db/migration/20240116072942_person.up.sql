CREATE TABLE IF NOT EXISTS people
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255),
    surname     VARCHAR(255),
    patronymic  VARCHAR(255),
    age         INTEGER,
    gender      VARCHAR(255),
    nationality VARCHAR(255)
);