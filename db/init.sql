CREATE TABLE users (
    id SERIAL,
    name VARCHAR(60),
    balance REAL
);

CREATE TABLE quests (
    id SERIAL,
    name VARCHAR(60),
    cost REAL
);

CREATE TABLE user_quest (
    user_id INT,
    quest_id INT
);
