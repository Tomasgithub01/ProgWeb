# ProgWeb

Trabajo práctico de cursada de la materia **Programación Web** de Ingeniería de Sistemas.  

**Integrantes:** Antúnez Monges Tomás, Buralli Agustín, Todesco Sofía Anabel.

---

## Instrucciones de uso

### Pre-requisitos

- Tener instalado un editor de texto, por ejemplo, [VSCode](https://code.visualstudio.com/download).  
- Tener instalado un compilador de [Go](https://go.dev/doc/install). Se recomienda la última versión.  
- Tener instalado [Git](https://git-scm.com/downloads) para clonar el repositorio.  
- Tener instalada la última versión de [Docker](https://docs.docker.com/engine/install/ubuntu/) y [Docker Compose](https://docs.docker.com/compose/install/linux/#install-using-the-repository). El proyecto fue testeado con **Docker Compose v2.24.5**.

---

### Pasos a seguir

1. **Clonar el repositorio**  

    Abrir una terminal en la carpeta deseada y ejecutar:  

    ```bash
    git clone https://github.com/Tomasgithub01/ProgWeb.git
    ```

2. **Abrir el proyecto en VSCode**  

    ```bash
    cd ProgWeb/
    code .
    ```

3. **Ejecutar el proyecto con Docker Compose**

    Si se quiere ejecutar con un solo comando utilizar:
   IMPORTANTE: Este comando posee un sleep arbitrario entre levantar el docker y realizar el testeo con hurl.

    ```bash
      make testdev
      ```   

    Este ejemplo contiene los comandos para ejecutar en modo de desarrollo (incluye air)
    - Construir e iniciar los contenedores :

      ```bash
      make dev
      ```

    - En una nueva terminal correr testeos con hurl
      ```bash
      make test
      ```

    - Abrir el navegador en:  
      ```
      http://localhost:8080
      ```

    - Para detener los contenedores después de su uso:

      Ctrl + C en la terminal que está corriendo y luego:
      
      ```bash
      make down
      ```

---

# 🎮 Funcionamiento General

- La app utiliza la **[API de Steam](https://steamcommunity.com/dev)** para buscar los juegos.  
  Desde allí hacemos una petición que nos devuelve un JSON del estilo:

  ```json
  {
    "id": 367520,
    "type": "game",
    "name": "Hollow Knight",
    "tiny_image": "https://cdn.akamai.steamstatic.com/steam/apps/367520/capsule_sm_120.jpg",
    "header_image": "https://cdn.akamai.steamstatic.com/steam/apps/367520/header.jpg",
    "hero_image": "https://cdn.akamai.steamstatic.com/steam/apps/367520/library_hero.jpg",
    "url": "https://store.steampowered.com/app/367520/Hollow_Knight/",
    "price": "14.99 USD"
  }
Las respuestas JSON se transforman al formato que entiende nuestro backend y se insertan en la base de datos a través del endpoint:


POST /game


Esto permite no tener que guardar todas las imágenes de los juegos en nuestro servidor de archivos, en detrimento de tener que buscarlas directamente desde Steam.
Por ejemplo, para las imágenes se utiliza la siguiente URL base:


https://cdn.cloudflare.steamstatic.com/steam/apps/${game.id}/hero_capsule.jpg
Completando el placeholder ${game.id} con el ID del juego obtenido desde la query inicial.

En el buscador de la parte principal, se puede buscar un juego y, al hacer clic, se envía al servidor Go, y en consecuencia, a la base de datos.

## Estructura del Proyecto

```text
ProgWeb/
├── db/                               # Esquemas, queries y código Go generado por SQLC
├── .air.toml                         # Archivo para automatizar los cambios en el código
├── .gitignore                        # Archivos de dependencias a ignorar
├── index.html                        # Página principal de la app
├── main.go                           # Código principal en Go (servidor HTTP)
├── go.mod                            # Dependencias de Go
├── Dockerfile                        # Imagen para la app en Go
├── Explicacion de la Aplicación.pdf  # Documento breve donde se detalla el funcionamiento de la app
├── docker-compose.yml                # Configuración de contenedores
├── Makefile                          # Comandos auxiliares
├── sqlc.yml                          # Archivo de configuración de SQLC
└── README.md                         # Documentación del proyecto
