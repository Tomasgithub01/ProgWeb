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
            <img src="${game.image}" alt="${game.name}">
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


const searchInput = document.getElementById("search");
const resultsDiv = document.getElementById("results");

let timeout;

searchInput.addEventListener("input", () => {
  clearTimeout(timeout);
  const query = searchInput.value.trim();
  if (!query) {
    resultsDiv.style.display = "none";
    return;
  }

  // Espera 300ms antes de buscar (para no saturar el servidor)
  timeout = setTimeout(async () => {
    const res = await fetch(`/steam/search?query=${encodeURIComponent(query)}`);
    const games = await res.json();

    resultsDiv.innerHTML = "";
    games.forEach(game => {
      const div = document.createElement("div");
      div.classList.add("result-item");
      div.innerHTML = `
        <img src="${game.tiny_image}" alt="${game.name}">
        <span>${game.name}</span>
      `;
      div.addEventListener("click", () => selectGame(game));
      resultsDiv.appendChild(div);
    });
    resultsDiv.style.display = "block";
  }, 300);
});

function selectGame(game) {
  resultsDiv.style.display = "none";
  searchInput.value = game.name;

  fetch("/games", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: game.name || "",              // nunca null
      description: "Imported from Steam",
      //image: game.tiny_image || "",       // usa string vacío si no hay imagen
      image: `https://cdn.cloudflare.steamstatic.com/steam/apps/${game.id}/hero_capsule.jpg` || "",
      link: `https://store.steampowered.com/app/${game.id}` || "" // también asegurado
    })
  }).then(res => {
    if (!res.ok) {
      res.text().then(t => console.error("Error inserting game:", t));
    }
  });
}