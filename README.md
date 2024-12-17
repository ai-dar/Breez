# Breez – Microblogging Platform

**Breez** is a lightweight microblogging platform inspired by Twitter. It allows users to post and retrieve short messages via a RESTful API and displays a basic frontend interface for interaction.

---

## 🌟 Project Overview

### Purpose
Breez demonstrates building a microblogging backend service with Go, RESTful APIs, database migrations, and a simple frontend.

### What It Does
- Supports `GET` and `POST` methods for messages.
- Stores data securely using a SQL database.
- Provides a minimal HTML/CSS interface for user interaction.

### Target Audience
- Developers learning Go (Golang).
- Students exploring backend development.
- Teams building scalable web applications.

---

## 👥 Team Members

- **Team Member 1**: Aidar Sabyrgali  
- **Team Member 2**: Miras Zhumaseitov 

---

## 🖥️ Screenshot of Main Page

![Main Page Screenshot](./static/main_page_screenshot.png)

---

## 🚀 How to Start the Project

### Prerequisites
1. Install **Go** (v1.16+): [Download Go](https://golang.org/dl/).
2. Install a SQL database (e.g., PostgreSQL or MySQL).
3. Install `golang-migrate` for database migrations.

---

### Steps to Run the Project

#### 1. Clone the Repository
    ```bash
       git clone https://github.com/your-team/breez.git
      cd breez
### 2. Set Up Environment Variables
### Create a .env file in the project root:

   ```
    DB_HOST=localhost
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=breez_db
    DB_PORT=5432
3. Run Database Migrations
Apply database migrations:
   ```
     bash
     migrate -path ./db/migrations -database "postgres://username:password@localhost:5432/breez_db?sslmode=disable" up
### 4. Start the Server
### Run the Go server:
   ```
    bash
    go run main.go
### 5. Access the Web Interface
### Open the following URL in your browser:
   ```
    arduino
    http://localhost:8080
## 📡 API Endpoints
### Method	Endpoint	Description
### GET	/json	Retrieve all messages.
### POST	/json	Add a new JSON message.
## 🛠️ Tools and Resources Used
### Go: Backend development.
### gorilla/mux: HTTP routing.
### SQL: Database for storing data.
### golang-migrate: Database migrations.
### HTML/CSS: Frontend for displaying content.

## 💡 Future Improvements
### Add user authentication (login/logout).
### Implement message persistence in a database.
### Create a dynamic frontend with JavaScript.

