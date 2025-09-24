# ProgWeb

Trabajo práctico de cursada de la materia **Programación Web** de Ingeniería de Sistemas.  

**Integrantes:** Antúnez Monges Tomás, Buralli Agustín, Todesco Sofía Anabel.

---

## Instrucciones de uso

### Pre-requisitos

- Tener instalado un editor de texto, por ejemplo, [VSCode](https://code.visualstudio.com/download).  
- Tener instalado un compilador de [Go](https://go.dev/doc/install). Se recomienda la última versión.  
- Tener instalado [Git](https://git-scm.com/downloads) para clonar el repositorio.  
- Tener instalada la última versión de Docker y Docker Compose. El proyecto fue testeado con **Docker Compose v2.24.5**.

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

    - Construir e iniciar los contenedores:

      ```bash
      docker compose up --build
      ```

    - Abrir el navegador en:  
      ```
      http://localhost:8080
      ```

    - Para detener los contenedores después de su uso:

      ```bash
      docker compose down
      ```

---

## Estructura del Proyecto

```text
ProgWeb/
├── db/                 # Esquemas, queries y código Go generado por SQLC
├── .air.toml           # Archivo para automatizar los cambios en el código
├── .gitignore          # Archivos de dependencias a ignorar
├── index.html          # Página principal de la app
├── main.go             # Código principal en Go (servidor HTTP)
├── go.mod              # Dependencias de Go
├── Dockerfile          # Imagen para la app en Go
├── docker-compose.yml  # Configuración de contenedores
├── Makefile            # Comandos auxiliares
├── sqlc.yml            # Archivo de configuración de SQLC
└── README.md           # Documentación del proyecto
