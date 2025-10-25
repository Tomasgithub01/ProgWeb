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
    const html = await response.text();


    //  Renderizamos dinámicamente los datos en el DOM
    container.innerHTML = html;
    // const list = document.createElement("ul");

    // data.forEach((item) => {
    //   const li = document.createElement("li");
    //   li.textContent = `${item.id} - ${item.nombre}`;
    //   list.appendChild(li);
    // });

    // container.appendChild(list);
  } catch (error) {
    console.error(error);
    container.textContent = "Error al cargar las entidades.";
  }
});