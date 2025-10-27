# Borrar lo que haya en la base
#DELETE http://localhost:8080/games/1
#HTTP 204

#DELETE http://localhost:8080/games/2
#HTTP 204

#DELETE http://localhost:8080/games/3
#HTTP 204

# Crear Juego 1
POST http://localhost:8080/games
Content-Type: application/json

{
  "name": "Minecraft",
  "description": "Juego de construcción con bloques y supervivencia.",
  "image": {"String":"https://example.com/minecraft.jpg","Valid":true},
  "link": {"String":"https://www.minecraft.net","Valid":true}
}

HTTP/1.1 201



# Crear Juego 2
POST http://localhost:8080/games
Content-Type: application/json

{
  "name": "Stardew Valley",
  "description": "Simulación de granja con exploración y aventuras.",
  "image": {"String":"https://example.com/stardew.jpg","Valid":true},
  "link": {"String":"https://www.stardewvalley.net","Valid":true}
}

HTTP/1.1 201


# Crear Juego 3
POST http://localhost:8080/games
Content-Type: application/json

{
  "name": "Hollow Knight",
  "description": "Metroidvania con exploración y combates desafiantes.",
  "image": {"String":"https://example.com/hollowknight.jpg","Valid":true},
  "link": {"String":"https://www.hollowknight.com","Valid":true}
}

HTTP/1.1 201



# Obtener todos los juegos
GET http://localhost:8080/games

HTTP/1.1 200


# Obtener el juego con id 3
GET http://localhost:8080/games/3
HTTP 200
[Asserts]
jsonpath "$.name" == "Hollow Knight"


# Update al juego con id 3
PUT http://localhost:8080/games/3
Content-Type: application/json
{
  "name": "Hollow Knight",
  "image": {"String":"https://i.blogs.es/0427e3/hollow-knight--0-/375_375.jpeg","Valid":true},
  "link": {"String":"https://www.hollowknight.com","Valid":true}
}
HTTP 200