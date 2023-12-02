# hexarchgo
A hexagonal Architecture implementation for Go Lang applications 

## Overview

This project implements a web-based API for a task management system using GoLang. It focuses on providing RESTful endpoints for creating, updating, and organizing tasks. The aim is to create a scalable and efficient task management solution for small to medium-sized teams.

## Features

- RESTful API for task CRUD operations
- User authentication and authorization
- Task prioritization and categorization
- Database persistence for task storage
- Scalable architecture for handling increased loads

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Configuration](#configuration)
- [Testing](#testing)
- [Contributing](#contributing)
- [Acknowledgments](#acknowledgments)
- [License](#license)

## Getting Started

To get started with this project:

1. Clone this repository:

    ```bash
    git clone https://github.com/your-username/project-name.git
    ```

2. Navigate to the project directory:

    ```bash
    cd project-name
    ```

3. Set up the database:

    - Create a PostgreSQL database named `tasks_db`.
    - Update the database connection details in `config/config.go`.

4. Run the application:

    ```bash
    go run main.go
    ```

## Usage

Once the application is running, the API endpoints can be accessed using tools like cURL or Postman. Here's an example of how to create a new task:

```bash
curl -X POST http://localhost:8080/ -d '{"title":"Task Title","description":"Task Description","priority":1}'
```


