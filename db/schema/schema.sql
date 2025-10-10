CREATE DOMAIN state AS varchar(20) 
CHECK (VALUE = 'Stateless' 
    OR VALUE = 'Started'
    OR VAlUE = 'Full Completed'
    OR VALUE = 'Abandoned');
    
CREATE DOMAIN rating AS integer
CHECK (VALUE > 0 AND VALUE <= 5);

CREATE TABLE game (
    id serial not null,
    name varchar(30) not null,
    description varchar(200) not null,
    image varchar(100),
    link varchar(100),
    CONSTRAINT pk_game PRIMARY KEY (id)
);

CREATE TABLE user (
    id serial not null,
    name varchar(30) not null,
    password varchar(255) not null,
    CONSTRAINT pk_user PRIMARY KEY (id)
);

CREATE TABLE plays (
    id_game int not null,
    id_user int not null,
    state state not null default 'Stateless',
    rating rating,
    CONSTRAINT pk_plays PRIMARY KEY (id_game, id_user)
);

ALTER TABLE plays
    ADD CONSTRAINT fk_plays_game 
    FOREIGN KEY (id_game) REFERENCES game (id);

ALTER TABLE plays
    ADD CONSTRAINT fk_plays_user
    FOREIGN KEY (id_user) REFERENCES user (id);