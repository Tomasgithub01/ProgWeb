# Creacion Usuario 1
POST http://localhost:8080/users
[Form]
username: 123
password: 123


# Creacion Usuario 2
POST http://localhost:8080/users
[Form]
username: 345
password: 345

# Logueo como 123
POST http://localhost:8080/login
[Form]
username: 123
password: 123

# Crear Juego 1
POST http://localhost:8080/games
[Form]
name: Minecraft
description: Juego de construcción con bloques y supervivencia.
image: https://image.api.playstation.com/vulcan/ap/rnd/202407/0401/670c294ded3baf4fa11068db2ec6758c63f7daeb266a35a1.png
link: https://www.minecraft.net
HTTP 303

# Crear Juego 2
POST http://localhost:8080/games
[Form]
name: Hollow Knight
description: Metroidvania con exploración y combates desafiantes.
image: /static/images/hollow_knight_hero.png
link: https://www.hollowknight.com
HTTP 303




# Obtener todos los juegos
GET http://localhost:8080/games

HTTP/1.1 200


# Obtener el juego con id 2
GET http://localhost:8080/games/2
HTTP 200


# Update al juego con id 2
PUT http://localhost:8080/games/2
[Form]
name: Hollow Knight
description: Metroidvania con exploración y combates desafiantes.
image: https://i.blogs.es/0427e3/hollow-knight--0-/375_375.jpeg
link: https://www.hollowknight.com
HTTP 303

#DELETE http://localhost:8080/games/3
#HTTP 204