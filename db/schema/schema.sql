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
    state state not null default 'Stateless',
    rating rating,
    link varchar(100),
    CONSTRAINT pk_game PRIMARY KEY (id)
);