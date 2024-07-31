# Tasklist App Backend


1. [Introduction](#introduction)
2. [Implemented features](#implemented-features)
3. [Development environment](#development-environment)
4. [Endpoints](#endpoints)
    - [/login](#login)
    - [/register](#register)
    - [/tasks/{userID}](#tasksuserid)
    - [/tasks/{userID}/{taskID}](#tasksuseridtaskid)
    - [/users/{userID}](#usersuserid)


## Introduction
A simple HTTP JSON API backend for a Todo web app which I made for a research seminar course in my school,  
and it is my first project written in Go. The web app didn't require anything too complicated, there needed to be a way to register new users and log in existing ones, and they need to be able to perform basic CRUD operations on the tasks they have access to. My side goal was to keep the app somewhat minimal when it comes to frameworks and other dependencies as an extensive standard library is listed as one of the standout features for the Go language. <br>

This project has probably gone as far as I want to take it. I might improve the documentation a bit, maybe create an admin account type with elevated privileges over the normal user type and finish the frontend, but overall I want to move  on to something more interesting/useful to me and hopefully to you too.

## Implemented features:

- CRUD operations for users and their tasks 
- JWT-based user Authentication  


## Development environment:
I ran the Go backend on the host, while a MySQL Docker container from the official image (https://hub.docker.com/_/mysql) served as the database server. The file 'dotenvBase.txt' has a field for every environment variable necessary for running the application. You just need to populate the fields and rename the file to '.env'.

## Endpoints
### /login
    
    Example: localhost:4200/login

    #### POST - Login
    Request Body example:
    {
        "email": "example@tasklist.com",
        "pasword": "Example1"
    }
    Response:
    {
        "username": "example",
        "token": {JWT-token}
    }        
### /register 

    Example: localhost:4200/register

    #### POST - Register a new user
    Request Body example:
    {
        "username": "ExampleUser"
        "email": "example@tasklist.com",
        "pasword": "Example1"
    }
    Response:
    {
        "username": "exampleUser",
        "email": "example@tasklist.com",
        "password": "Example1"
    }
### /tasks/{userID}    (JWT-Protected)
    Example: localhost:4200/tasks/1e2918cd-d27f-47e7-8318-cfd4d7056617


    #### GET - Get all tasks by users ID


    #### POST - Create a task for the user
    Request body example:
    {
        "title": "Math homework",
        "description": "page 51 assignments 1,2,3",
        "deadline": "2024-01-21" 
    }
    Response:
    {
        "task_id": "1eacc959-f665-4956-9303-1db47653abe0",
        "title": "Math homework",
        "description": "page 51 assignments 1,2,3",
        "deadline": "2024-01-21T00:00:00Z",
        "created_at": "2024-07-31T10:29:37Z",
        "updated_at": "2024-07-31T10:29:37Z",
        "user_id": "1e2918cd-d27f-47e7-8318-cfd4d7056617"
    }

### /tasks/{userID}/{taskID}     (JWT-Protected)

    
    Example: localhost:4200/tasks/1e2918cd-d27f-47e7-8318-cfd4d7056617/5f95a0f5-bd8b-4c2f-9973-f4b40fdb5404

    This endpoint looks a bit messy, since it has two UUID's in the URL. It's done this way because of how I implemented the JWT authentication.

    #### GET - Get users task selected by taskID

    #### PUT - Update task by taskID (Accepts partial objects)
    Request body example:
    {
        "description": "page 56 assignments 4,5,6",
    }
    Response (still has some room for improvement, currently returns the object in the state it was in before updating):
    {
        "task_id": "5f95a0f5-bd8b-4c2f-9973-f4b40fdb5404",
        "title": "Math homework",
        "description": "page 51 assignments 1,2,3",
        "deadline": "2024-12-01T00:00:00Z",
        "created_at": "2024-07-15T13:20:40Z",
        "updated_at": "2024-07-15T13:24:27Z",
        "user_id": "1e2918cd-d27f-47e7-8318-cfd4d7056617"
    }
 
    #### DELETE - Delete a task by taskID

### /users/{userID}      (JWT-Protected)

    Example: localhost:4200/users/1e2918cd-d27f-47e7-8318-cfd4d7056617

    #### GET - Get user by userID

    #### PUT - Update user information by userID (Accepts partial objects)
        Request Body example:
        {
            "username": "NewUserName"
        }
        Response:
        {
            "updated": "1e2918cd-d27f-47e7-8318-cfd4d7056617"
        }
    #### DELETE - Delete an user by userID

