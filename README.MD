# Go Web Server

## Project Overview
A **Go web server learning exercise** with room for improvement. It implements a simple RESTful API for user management with an in-memory database. To ensure a deep understanding of Go's syntax and standard library, WebServers, this project is being developed without the use of AI assistance.

## Architecture
- **Database**: In-memory cache (`userCache` map) storing User objects
- **Thread Safety**: Uses `sync.RWMutex` to ensure thread-safe access to shared data
- **Router**: Built with Go's `http.NewServeMux`

## API Endpoints

### Root Route
- **Method**: `GET`
- **Path**: `/`
- **Response**: "hello world"

### Create User
- **Method**: `POST`
- **Path**: `/users`
- **Request Body**: JSON with user name
- **Validation**: Name field is required
- **Response**: 204 No Content on success

### Get User
- **Method**: `GET`
- **Path**: `/users/{id}`
- **Response**: JSON representation of user
- **Status Codes**:
  - 200 OK - User found
  - 404 Not Found - User not found

### Delete User
- **Method**: `DELETE`
- **Path**: `/users/{id}`
- **Response**: 204 No Content on success
- **Status Codes**:
  - 204 No Content - User deleted
  - 400 Bad Request - User not found

## Server Configuration
- **Port**: 8080
- **Data Model**: User has a `name` field (string)

## Key Features
- Thread-safe read/write operations using `RWMutex`
- JSON serialization/deserialization
- Basic error handling with appropriate HTTP status codes
