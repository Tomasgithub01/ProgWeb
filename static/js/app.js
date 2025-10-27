document.addEventListener("DOMContentLoaded", async () => {
  const container = document.getElementById("gameList");

  try {
    //  Hacemos la petición a nuestro backend
    const response = await fetch("http://localhost:8080/games");

    //  Verificamos que la respuesta sea exitosa
    if (!response.ok) {
      throw new Error("Error en la respuesta del servidor");
    }

    //  Interpretamos la respuesta como JSON
    //const html = await response.text();
    const data = await response.json();
    console.log(data);
    

    container.innerHTML = ""; // clear previous content

    if (data.length === 0) {
      // If there are no elements, show a message
      const p = document.createElement("p");
      p.textContent = "There are no elements to display.";
      container.appendChild(p);
    } else {
      // If there are elements, render the list
      const list = document.createElement("ul");
      data.forEach((item) => {
        const li = document.createElement("li");
        li.textContent = `${item.id} - ${item.name}`;
        list.appendChild(li);
      });
      container.appendChild(list);
    }
  } catch (error) {
    console.error(error);
    container.textContent = "Error al cargar las entidades.";
  }
});