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
  "image": "https://image.api.playstation.com/vulcan/ap/rnd/202407/0401/670c294ded3baf4fa11068db2ec6758c63f7daeb266a35a1.png",
  "link": "https://www.minecraft.net"
}

HTTP/1.1 201



# Crear Juego 2
POST http://localhost:8080/games
Content-Type: application/json

{
  "name": "Stardew Valley",
  "description": "Simulación de granja con exploración y aventuras.",
  "image": "https://cdn.cloudflare.steamstatic.com/steam/apps/413150/hero_capsule.jpg",
  "link": "https://www.stardewvalley.net"
}

HTTP/1.1 201


# Crear Juego 3
POST http://localhost:8080/games
Content-Type: application/json

{
  "name": "Hollow Knight",
  "description": "Metroidvania con exploración y combates desafiantes.",
  "image": "https://example.com/hollowknight.jpg",
  "link": "https://www.hollowknight.com"
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
  "image": "https://i.blogs.es/0427e3/hollow-knight--0-/375_375.jpeg",
  "link": "https://www.hollowknight.com"
}
HTTP 200