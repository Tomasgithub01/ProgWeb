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