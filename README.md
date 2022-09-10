# rakamin-mini-project
Rakamin Mini Project Backend

[Live Server](https://rakamin-mini-project.portalnesia.com)

## Installation

1. Clone Repository

    ```bash
    git clone https://github.com/putuadityabayu/rakamin-mini-project.git
    ```

2. Install Mysql >= 8.0

    You can use [sample data](#database-dump) for easy development (optional)

3. Add `.env` file

    [Example env file](.env.example)

4. Run mod tidy
  
    ```bash
    go mod tidy
    ```

5. Run Server

    ```bash
    go run main.go
    ```

6. Server running on ***localhost:$PORT***

-----

## Database Dump

You can use sample data to import in mysql database. SQL data can be seen in [this file](database.sql)

### Sample User

1. Name: User 1   
    Username: user1   
    Password: user1

2. Name: User 2  
    Username: user2   
    Password: user2

3. Name: User 3   
    Username: user3   
    Password: user3

------

## Case Study

You can see the response of this case study using postman [here](https://postman.com/portalnesia/workspace/rakamin-mini-projects/documentation/13670841-65bbdf05-492e-44c3-b1db-f4d979fbaa25).

*Noted: This case study uses sample data contained in this repository*

### User Story #1

As a user, I want to be able to send message to other user, so that I will be able to share information with others.

1. Login

    Send POST request to `/login` endpoint with User data above.

    Example request body:
    ```json
    {
      "username":"user1",
      "password": "user1"
    }
    ```

    Example response:
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJSYWthbWluIE1pbmkgUHJvamVjdCIsInN1YiI6IjEiLCJhdWQiOlsiaHR0cHM6Ly9yYWthbWluLmNvbSJdLCJleHAiOjE2NjI2NjAyMzh9.F9kzLO4Mn8PA8WHJC4IHRcMKfha7_49zvrIuicWvzsE",
      "user": {
        "id": 1,
        "name": "User 1",
        "username": "user1"
      }
    }
    ```

2. Authenticated

    All access to the `/conversation/**` endpoint must include an authorization header:

    ```http
    Authorization: Bearer JWT_TOKEN
    ```

3. Send `POST` request to `/conversation`

    To send new messages, send request with json body:

    ```json
    {
      "user_id": 2,
      "message": "Example message"
    }
    ```

    - Scenario 1

      Send the message with fill message.

      Example response:
      ```json
      {
        "id": 30,
        "read_status": false,
        "timestamp": "2022-09-09T14:01:04.932+07:00",
        "message": "Example message",
        "sender": {
          "id": 1,
          "name": "User 1",
          "username": "user1"
        },
        "conversation": {
          "id": 11,
          "unread": 0,
          "users": [
            {
              "id": 1,
              "name": "User 1",
              "username": "user1"
            },
            {
              "id": 2,
              "name": "User 2",
              "username": "user2"
            }
          ]
        }
      }
      ```

    - Scenario 2

      Send the message without fill message.

      Example response:
      ```json
      {
        "error": "Message cannot be empty"
      }
      ```

### User Story #2

As a user, I want to be able to reply message in existing conversation, so that I will be able to respond previous message.

1. Login with user2 account
2. Authenticated
3. Send message (reply) to conversation

    Send `POST` request to `/conversation/:id`, where `id` is conversation ID.

    Example request body:
    ```json
    {
      "message":"Example reply message"
    }
    ```

    Example response:
    ```json
    {
      "id": 31,
      "read_status": false,
      "timestamp": "2022-09-09T15:36:49.448+07:00",
      "message": "Tes reply to user 1 from user 2",
      "sender": {
        "id": 2,
        "name": "User 2",
        "username": "user2"
      },
      "conversation": {
        "id": 11,
        "unread": 0,
        "users": [
          {
            "id": 1,
            "name": "User 1",
            "username": "user1"
          },
          {
            "id": 2,
            "name": "User 2",
            "username": "user2"
          }
        ]
      }
    }
    ```


### User Story #3

As a user, I want to be able to list messages from specific user, so that I will be able to read our conversation.

1. Login
2. Authenticated
3. List all messages in specific conversation

    Send `GET` request to `/conversation/:id`, where `id` is conversation ID.

    Optional query:

    - `page`: Page you want to request. Default 1
    - `page_size`: Size of data in one call requests. Default 15


    Example response:
    ```js
    {
      "page": 1, // Requested page
      "page_size": 15, // Size of data in one call requests
      "total": 2, // Total data
      "total_page": 1, // Maximum pages that can be requested. Page requests larger than this value, will result in empty data
      "data": [
        {
          "id": 30,
          "read_status": true,
          "timestamp": "2022-09-09T14:01:04.932+07:00",
          "message": "tes user 2 msg 2",
          "sender": {
            "id": 1,
            "name": "User 1",
            "username": "user1"
          }
        },
        ...messages
      ]
    }
    ```

### User Story #4

As a user, I want to be able to list conversations where I involved, so that I will be able to search or find user to chat with.

1. Login
2. Authenticated
3. List all conversation

    Send `GET` request to `/conversation`.

    Optional query:

    - `page`: Page you want to request. Default 1
    - `page_size`: Size of data in one call requests. Default 15

    Example response:
    ```js
    // authorization_user is the user associated with the authorization header token
    {
      "page": 1, // Requested page
      "page_size": 15, // Size of data in one call requests
      "total": 2, // Total data
      "total_page": 1, // Maximum pages that can be requested. Page requests larger than this value, will result in empty data
      "data": [
        {
          "id": 11,
          "unread": 0, // Unread count. Number of messages (in this conversation) that have not been read by authorization_user
          "users": [ // Conversation participants information
            {
              "id": 1,
              "name": "User 1",
              "username": "user1"
            },
            {
              "id": 2,
              "name": "User 2",
              "username": "user2"
            }
          ],
          "message": {
            "id": 31,
            "read_status": true, // Status whether the message has been read or not read by the recipient 
            "timestamp": "2022-09-09T15:36:49.448+07:00",
            "message": "Tes reply to user 1 from user 2",
            "sender": {
              "id": 2,
              "name": "User 2",
              "username": "user2"
            }
          }
        },
        ...conversations
      ]
    }
    ```

------

## TODO

- [x] Added CI
- [x] Added CD