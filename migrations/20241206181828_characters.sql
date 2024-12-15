-- +goose Up
-- +goose StatementBegin

CREATE TABLE actors (
    actor_id SERIAL PRIMARY KEY,
    actor_name VARCHAR(255) UNIQUE,
    actor_link TEXT
);
CREATE TABLE characters (
    character_id SERIAL PRIMARY KEY,
    character_name VARCHAR(255) NOT NULL UNIQUE,
    house_name VARCHAR(255),
    character_image_thumb TEXT,
    character_image_full TEXT,
    character_link TEXT,
    nickname VARCHAR(255),
    royal BOOLEAN
);

CREATE TABLE relationships (
    relationship_id SERIAL PRIMARY KEY,
    character_id INT,
    character_relationship_id INT,
    relationship_type VARCHAR(255),
    FOREIGN KEY (character_id) REFERENCES characters(character_id) ON DELETE CASCADE,
    FOREIGN KEY (character_relationship_id) REFERENCES characters(character_id) ON DELETE CASCADE
);

CREATE INDEX characters_character_name_idx ON characters(character_name);

CREATE TABLE characters_actors (
    character_actor_id SERIAL PRIMARY KEY,
    character_id INT,
    actor_id INT,
    FOREIGN KEY (character_id) REFERENCES characters(character_id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors(actor_id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE characters CASCADE;
DROP TABLE relationships CASCADE;
DROP TABLE characters_actors CASCADE;
DROP TABLE actors CASCADE;
-- +goose StatementEnd
