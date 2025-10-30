/*
WEBLOAD Initial data
Sends request to backend to load existing games and display them
*/
document.addEventListener("DOMContentLoaded", async () => {
  const gameList = document.getElementById("gameList");

  try {
  // Api call to our backend
    const response = await fetch("http://localhost:8080/games");

    
    if (!response.ok) {
      throw new Error("Server response error");
    }

    
    const games = await response.json();
    console.log(games);
    

    gameList.innerHTML = ""; 

    /*
    Adds delete button to each card,  onclick calls delete game function with the game id and the card element as parameters to remove it
    */

    if (games.length === 0) {
      const p = document.createElement("p");
      p.textContent = "There are no elements to display.";
      gameList.appendChild(p);
    } else {
      gameList.innerHTML = games.map(game => `
        <div class="game-card" data-id="${game.id}">
          <div class="game-image">
            <img src="${game.image}" alt="${game.name}">
          </div>
          <div class="game-name-btn">
            <div class="game-title">${game.name}</div>
            <div class="dropend">
                <button type="button" class="btn btn-outline-primary" style="--bs-btn-border-color: none;" data-bs-toggle="dropdown" aria-expanded="false">
                    <i class="fas fa-ellipsis-v"></i>
                </button>
                <ul class="dropdown-menu">
                    <li><a class="dropdown-item" onclick="deleteGame(${game.id}, this.closest('.game-card'))">Delete</a></li>    
                </ul>
            </div>
          </div>
        </div>
      `).join('');
    }
  
  } catch (error) {
    console.error(error);
    gameList.textContent = "Error loading games.";
  }
});


const searchInput = document.getElementById("search");
const resultsDiv = document.getElementById("results");

let timeout;

/*
SEARCH GAME JS FUNCTION
Sends search request to STEAM api endpoint and displays results
Adds an event listener to the search input field that triggers on input.
*/
searchInput.addEventListener("input", () => {
  clearTimeout(timeout);
  const query = searchInput.value.trim();
  if (!query) {
    resultsDiv.style.display = "none";
    return;
  }

  /*
  Shows results in form of card divs
  Add event listener to call select game function on click so to insert it
  */
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

/*
INSERT GAME JS FUNCTION
Sends post request to api endpoint to insert a new game
We use steam api data to fill in the fields previuously searched
*/
function selectGame(game) {
  resultsDiv.style.display = "none";
  searchInput.value = game.name;

  fetch("/games", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      name: game.name || "",              
      description: "Imported from Steam",
      image: `https://cdn.cloudflare.steamstatic.com/steam/apps/${game.id}/hero_capsule.jpg` || "",
      link: `https://store.steampowered.com/app/${game.id}` || "" 
    })
  }).then(res => {
    if (!res.ok) {
      res.text().then(t => console.error("Error inserting game:", t));
    } else {
      location.reload(); 
    }
  });
}
/*
DELETE GAME JS FUNCTION
Sends delete request to api endpoint and removes card from DOM on success
(could also be done by reloading the DOM instead of deleting part of it)
*/
async function deleteGame(id, cardEl) {
  if (!confirm('¿Delete this game?')) return;
  try {
    const res = await fetch(`/games/${id}`, { method: 'DELETE' });
    if (!res.ok) {
      const text = await res.text();
      console.error('Error deleting game:', text); 
      alert('Error deleting game'); 
      return;
    }
    if (cardEl && cardEl.remove) cardEl.remove();
  } catch (err) {
    console.error('Error deleting game:', err);
    alert('Error deleting game');
  }
}