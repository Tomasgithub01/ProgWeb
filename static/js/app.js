document.addEventListener("DOMContentLoaded", async () => {
  const gameList = document.getElementById("gameList");

  try {
    //  Hacemos la petición a nuestro backend
    const response = await fetch("http://localhost:8080/games");

    //  Verificamos que la respuesta sea exitosa
    if (!response.ok) {
      throw new Error("Error en la respuesta del servidor");
    }

    //  Interpretamos la respuesta como JSON
    //const html = await response.text();
    const games = await response.json();
    console.log(games);
    

    gameList.innerHTML = ""; // clear previous content

    if (games.length === 0) {
      // If there are no elements, show a message
      const p = document.createElement("p");
      p.textContent = "There are no elements to display.";
      gameList.appendChild(p);
    } else {
      gameList.innerHTML = games.map(game => `
        <div class="game-card">
          <div class="game-image">
            <img src="https://i.blogs.es/0427e3/hollow-knight--0-/375_375.jpeg" alt="${game.name}">
          </div>
          <div class="game-title">${game.name}</div>
        </div>
      `).join('');
    }
  
  } catch (error) {
    console.error(error);
    gameList.textContent = "Error al cargar las entidades.";
  }
});