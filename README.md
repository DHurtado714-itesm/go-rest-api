# 📚 Book Management API with Go and MariaDB

Welcome to the **Book Management API**! This is a simple RESTful API built with [Go](https://golang.org/) and [Gin](https://github.com/gin-gonic/gin) that allows you to manage a collection of books using a MariaDB database. 🚀

## 🛠️ Features

- 📖 **Get All Books**: Retrieve a list of all books in the collection.
- ➕ **Add a New Book**: Add a new book to the collection.
- ✏️ **Edit a Book**: Update the details of an existing book.
- ❌ **Delete a Book**: Remove a book from the collection.

## 🚀 Getting Started

### Prerequisites

Before you begin, make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.18 or later)
- [Docker](https://www.docker.com/get-started) (to run MariaDB)
- [Git](https://git-scm.com/)

### 🔧 Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/your-username/book-management-api.git
   cd book-management-api
    ```

2. **Start MariaDB with Docker**

    ```bash
    docker run --name book-db -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=bookdb -p 3306:3306 -d mariadb:latest
    ```

3. **Run the Application**

    ```bash
    go run main.go
    ```

4. **Access the API**
  
The API will be available at `http://localhost:8080`.

## 📚 API Endpoints

The following endpoints are available:

- **GET /books**: Retrieve a list of all books.
- **POST /books**: Add a new book.
- **PUT /books**: Update a book by ID.
- **DELETE /books/:id**: Delete a book by ID.
