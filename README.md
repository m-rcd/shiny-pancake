# Notes

- [Requirements](#requirements)
- [Technologies](#technologies)
- [Usage](#usage)
- [Testing](#testing)
- [Approach](#approach)

REST API backend application that can be used to manage personal notes.

## Requirements

The API have the following key features:

- Save a new note
- Update a previously saved note
- Delete a saved note
- Archive a note
- Unarchive a previously archived note
- List saved notes that aren't archived
- List notes that are archived

## Technologies

- [Golang](https://go.dev/). Even though I have more experience with Ruby, I decided to implement the task in Golang instead because I thought it would be more challenging and a good opportunity to learn.
- [Gorilla mux](https://pkg.go.dev/github.com/gorilla/mux#section-readme) to handle requests. I chose this one because it is widely used and well supported.
- [ginkgo](https://github.com/onsi/ginkgo) for testing. I have opted for this one rather then the inbuilt go test because it allows for more descriptive tests. 
- [go-sql-mysql](https://github.com/go-sql-driver/mysql) I chose this one because it is well maintained and supporrted.
- [go-sqlmock](github.com/DATA-DOG/go-sqlmock) to mock sql queries in unit tests.
- [counterfeiter](github.com/maxbrunsfeld/counterfeiter/) to generate a fake database interface for unit tests.


## Usage

1. Clone the repo
    ```shell
    git clone git@github.com:m-rcd/notes.git
    cd notes
    ```

1. Build the app

    ```shell
    make build
    ```

1. Start the server

    ```shell
    ./notes
    ```
    The server will listen on port `10000`. 

    The server can take flags:
    - `--db` which can be `local` or `sql`. If not specified, the notes would be stored locally by default. 
    -  `--directory` to allow user to save notes in a specified location. If not specified, the notes would be saved in the default location `/tmp`. This flag is only used in the case of local storage.

    To save in a different directory: 
    ```shell
    ./notes --db local --directory <dir name>
    ```

    To use `sql` as database: 
    ```shell
    DB_USERNAME=<username> DB_PASSWORD=<password> ./notes --db sql
    ```

    To use `sql` a database `notes` needs to be created before the server has started. The table `notes` will be created as part of the app.

1. Create a note

    Open a new terminal and run the following command: 

    ```shell
    curl -X POST -H "Content-Type: application/json" -d '{"name":"note1","content":"I am a note!","user":{"username":"Sabriel"}}' http://localhost:10000/note
    ```

    The POST request will return a JSON response: 
    ```json
    {
        "type":"success",
        "StatusCode":200,
        "data":[
            {
                "id":"4ac82864-0354-43af-5582-fc721dfc4cf4",
                "name":"note1",
                "content":"I am a note!",
                "user":{
                    "username":"Sabriel"
                    },
                "archived":false
            }],
        "message":"The note was successfully created"
    }
    ```


    **Local Storage**

    This will create a new file with a name **name_id.txt** in directory `/tmp/note/username/active/`.

    In this case,  `note1_4ac82864-0354-43af-5582-fc721dfc4cf4.txt` in `/tmp/notes/Sabriel/active/` folder.
    The file will contain the content specified in the body of the request ("I am a useful note!").

    **SQL**

    This will create a record in the notes table.

    If the name was note sent:
    ```shell
    curl -X POST -H "Content-Type: application/json" -d '{"name":"","content":"I am a note!","user":{"username":"Sabriel"}}' http://localhost:10000/note
    ```
    The POST request will return a JSON response: 

    ```json
    {
        "type":"failed",
        "StatusCode":500,
        "data":[],
        "message":"name must be set"
    }
    ```
    The same validation is present for the user attribute.

1. Update a previously saved note
    ```
    curl -X PATCH -H "Content-Type: application/json" -d '{"name":"note1","content":"I am updated!","user":{"username":"Sabriel"}}' http://localhost:10000/note/4ac82864-0354-43af-5582-fc721dfc4cf4
    ```

    The PATCH request will return a JSON response: 
    ```json
    {
        "type":"success",
        "StatusCode":200,
        "data":[
            {
                "id":"4ac82864-0354-43af-5582-fc721dfc4cf4",
                "name":"note1",
                "content":"I am updated!",
                "user":{
                    "username":"Sabriel"
                    },
                "archived":false}],
        "message":"The note was successfully updated"
    }
    ```

    **Local Storage System**

    The content of `/tmp/notes/Sabriel/active/note1_4ac82864-0354-43af-5582-fc721dfc4cf4.txt` will be updated to "I am updated!".

    **SQL**

    The record will be updated in the database with the new data sent in the request.

1. Delete a saved note

    ```shell
    curl -X DELETE -H "Content-Type: application/json" -d '{"username":"Sabriel"}' http://localhost:10000/note/4ac82864-0354-43af-5582-fc721dfc4cf4
    ```

    The DELETE request will return a JSON response: 
    ```json
    {
        "type":"success",
        "StatusCode":200,
        "data":[],
        "message":"The note was successfully deleted"
    }
    ```
    **Local Storage System**

    The file `/tmp/notes/Sabriel/active/note1_4ac82864-0354-43af-5582-fc721dfc4cf4.txt` will be deleted.

    **SQL**

    The record will be deleted from database.

1. Archive a note 

    To archive a note, a PATCH request is used to move the note from `/tmp/notes/Sabriel/active/` to `/tmp/notes/Sabriel/archived/`. 
    Note that any other attribute sent in the body of the request will be ignored if `archived` is set to `true`.

    ```shell
    curl -X PATCH -H "Content-Type: application/json" -d '{"archived":true,"user":{"username":"Sabriel"}}' http://localhost:10000/note/4ac82864-0354-43af-5582-fc721dfc4cf4
    ```

    The PATCH request will return a JSON response: 
    ```json
    {
        "type":"success",
        "StatusCode":200,
        "data":[
            {"id":"4ac82864-0354-43af-5582-fc721dfc4cf4",
            "name":"note1",
            "content":"I am updated!",
            "user":{"username":"Sabriel"},
            "archived":true}],
        "message":"The note was successfully updated"
    }
    ```
    **SQL**

    The note will have the attribute `archived` set to true.

1. Unarchive a note 

    To unarchive a note, a PATCH request is used to move the note from `/tmp/notes/Sabriel/archived/` to `/tmp/notes/Sabriel/active/`. 

    Note that any other attribute sent in the body of the request will be ignored if `archived` is set to `false` and the note was archived.

    ```shell
    curl -X PATCH -H "Content-Type: application/json" -d '{"archived":false,"user":{"username":"Sabriel"}}' http://localhost:10000/note/4ac82864-0354-43af-5582-fc721dfc4cf4
    ```
    The PATCH request will return a JSON response: 
    ```json
    {
        "type":"success",
        "StatusCode":200,
        "data":[
            {
                "id":"4ac82864-0354-43af-5582-fc721dfc4cf4",
                "name":"note1",
                "content":"I am updated!",
                "user":{
                    "username":"Sabriel"
                    },
                "archived":false}],
        "message":"The note was successfully updated"
    }
    ```

     **SQL**

    The note will have the attribute `archived` set to false.


1. List saved notes that aren't archived

    ```shell
    curl -X GET -H "Content-Type: application/json" -d '{"username":"Sabriel"}' http://localhost:10000/notes/active
    ```

    The GET request will return a JSON response: 

    ```json
    [
        {
            "id":"4ac82864-0354-43af-5582-fc721dfc4cf4",
            "name":"note1",
            "content":"I am a useful note!",
            "user":{
                "username":"Sabriel"
                },
            "archived":false
        },
        {
            "id":"08c42190-19d6-4975-741b-ab87edfb9dc0",
            "name":"note2",
            "content":"I am another note!",
            "user":{
                "username":"Sabriel"
                },
            "archived":false
        }
    ]
    ```
    SQL storage not supported for this action yet.


1. List saved notes that are archived

    ```shell
    curl -X GET -H "Content-Type: application/json" -d '{"username":"Sabriel"}' http://localhost:10000/notes/archived
    ```

    The GET request will return a JSON response: 

    ```json
    [
        {
            "id":"4ac82864-0354-43af-5582-fc721dfc4cf4",
            "name":"note1",
            "content":"I am a useful note!",
            "user":{
                "username":"Sabriel"
                },
            "archived":true
        },
        {
            "id":"08c42190-19d6-4975-741b-ab87edfb9dc0",
            "name":"note2",
            "content":"I am another note!",
            "user":{
                "username":"Sabriel"
                },
            "archived":true
        }
    ]
    ```
    SQL storage not supported for this action yet.


## Testing

To run all tests: 

`DB_USERNAME` and `DB_PASSWORD` should be set in an `.env` file.

```shell
make test
```

## Approach

- I have opted to initially complete the API using a local storage instead of a SQL database for two reasons. 
One, it will allow a developer to use the app without having to setup the database before hand. Two, it allowed me to focus on the structure of my code. However, I implemented it using an interface so that once I have the structure completed, I can easily use that interface to plug in a database.

## Todo

- [x] Add Makefile
- [ ] Test for SQL 
- [ ] Add SQL support for all request
- [ ] Refactor integration test to be table test once SQL is working
- [ ] Add logging to help debugging in the case of server errors
- [ ] Add graceful shutdown
- [ ] User Auth
- [ ] TLS
- [ ] SQL: Add a user table and `user_id` column instead of `username` in notes table.
- [ ] Add routes to create a user
