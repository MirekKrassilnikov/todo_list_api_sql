
# ToDo List API

ToDo List API is a RESTful API for managing tasks. It allows users to create, read, update, and delete tasks, with support for recurring tasks.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MirekKrassilnikov/todo_list_api_sql.git
Navigate to the project directory:
bash
cd todo_list_api_sql
Make sure you have Go installed (version 1.15 and above). If Go is not installed, follow the official Go installation guide.
Install dependencies:
bash
go mod tidy
Create and configure the database, then make any necessary changes to the configuration files.
Running the Application

To run the server, execute the command:
bash
go run main.go
By default, the server will listen on port 8080.
Usage

The API provides the following endpoints:
Create a Task
POST /tasks
Request Body:
json
{
  "date": "2023-10-25",
  "title": "Task Title",
  "comment": "Comment for the task",
  "repeat": "d"
}
Get All Tasks
GET /tasks
Returns a list of all tasks.
Get a Task by ID
GET /tasks?id=1
Returns the task with the specified ID.
Update a Task
PUT /tasks
Request Body:
json
{
  "id": "1",
  "date": "2023-10-26",
  "title": "Updated Task Title",
  "comment": "Updated Comment",
  "repeat": "y"
}
Delete a Task
DELETE /tasks?id=1
Deletes the task with the specified ID.
Example Requests

Example request to create a task:
bash
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"date":"2023-10-25","title":"Task Title","comment":"Comment for the task","repeat":"d"}'
License

This project is licensed under the MIT License. For more details, see the LICENSE file.
Contribution

If you would like to contribute to the project, please create a pull request or open an issue with your suggestions.
Author

Miroslav Krasilnikov
