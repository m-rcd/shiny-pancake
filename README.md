# Notes

[Requirements](#requirements) | [Technologies](#technologies) | [Usage](#usage) | [testing](#testing)   |[approach](#approach)

REST API backend application that can be used to manage personal notes.

## <a name="requirements">**Requirements**</a>
---

The API have the following key features:

- Save a new note
- Update a previously saved note
- Delete a saved note
- Archive a note
- Unarchive a previously archived note
- List saved notes that aren't archived
- List notes that are archived

## <a name="technologies">**Technologies**</a> 
---

- [Gorilla mux](https://pkg.go.dev/github.com/gorilla/mux#section-readme) to handle requests.
- [net/http](https://pkg.go.dev/net/http) to handle HTTP client and server.
- [ginkgo](https://github.com/onsi/ginkgo) for testing. I have opted for this one rather then the inbuilt go test because it allows for more descriptive tests. 


## <a name="usage">**Usage**</a> 
---

1. Clone the repo
```
git clone git@github.com:m-rcd/notes.git
cd notes
```
2. Start the server

```
./notes
```
The server will listen on port `10000`. 

The command can take a flag `--directory` to allow user to save notes in a specified location. If not specified, the notes would be saved in the default location `/tmp`. 

```shell
./notes --directory <dir name>
```

3. Create a note

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

This will create a new file with a name **name_id.txt** in directory `tmp/note/username/active/`.

In this case,  `note1_4ac82864-0354-43af-5582-fc721dfc4cf4.txt` in `/tmp/notes/Sabriel/active/` folder.
The file will contain the content specified in the body of the request ("I am a useful note!").

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

4. Update a previously saved note
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

5. Delete a saved note

```
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

6. Archive a note 

To archive a note, a PATCH request is used to move the note from `/tmp/notes/Sabriel/active/` to `/tmp/notes/Sabriel/archived/`. 
Note that any other attribute sent in the body of the request will be ignored if `archived` is set to `true`.

```
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

6. Unarchive a note 

To unarchive a note, a PATCH request is used to move the note from `/tmp/notes/Sabriel/archived/` to `/tmp/notes/Sabriel/active/`. 

Note that any other attribute sent in the body of the request will be ignored if `archived` is set to `false` and the note was archived.

```
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

7. List saved notes that aren't archived

```
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

8. List saved notes that are archived

```
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

## <a name="testing">**Testing**</a>  

To run integration tests: 
```
ginkgo -r 
```

## <a name="approach">**Approach**</a>
---

- I have opted to initially complete the API using a local storage instead of a SQL database for two reasons. 
One, it will allow a developer to use the app without having to setup the database before hand. Two, it allowed me to focus on the structure of my code. However, I implemented it using an interface so that once I have the structure completed, I can easily use that interface to plug in a database.

