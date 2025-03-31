# Series backlog API
This is an API built with **Go**, using a Chi router and **SQLite** database. It provides endpoints for managing TV series, allowing CRUD operations.

## Getting started
The repository provides a Dockerfile with the necessary instructions to build and run the API inside a container.

**1. Clone the repository**.
```sh
git clone https://github.com/vicperezch/lab6-backend.git
cd lab6-backend
```

**2. Build the image**.
With Docker installed, run:
```sh
docker -t build <name> .
```

**3. Run the container**.
```sh
docker run -p 8080:8080 <name>
```

## API Endpoints
| Method | Endpoint |
|--------|----------|
| `GET`  | `/series` |
| `GET`  | `/series/{id}` |
| `POST` | `/series` |
| `PUT`  | `/series/{id}` |
| `DELETE` | `/series/{id}` |
| `PATCH` | `/series/{id}/upvote` |
| `PATCH` | `/series/{id}/downvote` |
| `PATCH` | `/series/{id}/status` |
| `PATCH` | `/series/{id}/episode` |

## Swagger Documentation
This API includes **Swagger UI** for easy documentation and testing.

Once the server is running, open in your browser:
```
http://localhost:8080/swagger/index.html
```

## Frontend example
![image](https://github.com/user-attachments/assets/91767474-3986-44b4-82b1-3c80ea91a077)
