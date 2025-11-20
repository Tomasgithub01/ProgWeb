# ProgWeb

Trabajo práctico de cursada de la materia **Programación Web** de Ingeniería de Sistemas.  

**Integrantes:** Antúnez Monges Tomás, Buralli Agustín, Todesco Sofía Anabel.

---

## 🧾 Instrucciones de uso

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


## 🔑 Primer Acceso y Autenticación

### Crear una Cuenta
- Haz clic en **Create one** para abrir el modal de registro.  
- Ingresa usuario y contraseña.  
- El sistema almacena tus credenciales en la base de datos.

### Iniciar Sesión
- Completa los campos con tus credenciales.  
- Haz clic en **Log In**.  
- Se envía una solicitud **POST /login** y, si el usuario existe, serás redirigido al dashboard principal (**/dashboard**).

---

## 🕹️ Funcionalidades del Dashboard

### Buscar y Agregar Juegos de Steam
- Usa la barra superior central para buscar juegos.  
- La aplicación consulta la API de Steam en tiempo real.  
- Al seleccionar un resultado, el juego se agrega automáticamente a tu catálogo personal.

### Agregar Juegos Personalizados
- En **My Games**, haz clic en la tarjeta **ADD NEW GAME**.  
- Se abre un modal para cargar datos de juegos no presentes en Steam.  
- La información se guarda mediante una solicitud **POST** al servidor.

### Modificar y Eliminar Juegos Propios
- Cada tarjeta de juego incluye un menú de tres puntos (⋮).  
- Desde el dropdown puedes **Modificar** o **Eliminar** el juego.

### Gestionar el Estado de los Juegos
Cada tarjeta posee un ícono de estado en la esquina superior derecha. Permite alternar entre:

- **Stateless** (representado por un bookmark)  
- **Started** (representado por un gamepad)  
- **Completed** (representado por una estrella)
- **Full Completed** (representado con una medalla)  
- **Abandoned** (representado por un corazón partido)

### Cerrar Sesión
- En la esquina superior derecha encontrarás el botón **Logout**.  
- Cierra la sesión y redirige a la pantalla principal.


## 🌳 Estructura del Proyecto

```text
ProgWeb/
├── db/                               # Esquemas (schema.sql), queries (queries.sql) y código Go generado por SQLC.
│   ├── queries
│   │   └── queries.sql               # Consultas SQL para obtener juegos y estados (plays).
│   ├── schema
│   │   └── schema.sql                # Definición de tablas.
│   └── sqlc                          # Código Go generado automáticamente.
├── logic/                            # Lógica de negocio específica (ej: gestión de APIs).
│   └── games.go                      # Funciones para interactuar con la API de Steam.
├── static/                           # Recursos estáticos: CSS, imágenes, archivos JS de frontend.
│   ├── css
│   ├── images                        # Incluye íconos de estado y elementos de la UI.
│   ├── js
│   ├── index.html                    # Layouts principales de la aplicación.
│   └── login.html
├── views/                            # Plantillas de la interfaz de usuario escritas en Templ.
│   ├── games.templ
│   └── layout.templ                  # Define la estructura base del dashboard.
├── main.go                           # Servidor HTTP principal, manejo de rutas y lógica de controlador.
└── docker-compose.yml                # Configuración de contenedores (Go App + PostgreSQL).

