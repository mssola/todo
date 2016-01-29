# ToDo [![Build Status](https://travis-ci.org/mssola/todo.svg?branch=master)](https://travis-ci.org/mssola/todo) [![GoDoc](https://godoc.org/github.com/mssola/todo?status.png)](http://godoc.org/github.com/mssola/todo)

## About this application

This application has a simple job: handle my "To Do list". This might seem
trivial but I got to a point where my To Do list was a scrambled mess that was
scattered throughout my system. To fix this situation I've built a service that
has the following specifications:

* There only exists one user in the database.
* The "To Do" list is organized in topics. Topics can be created, deleted,
  renamed, updated, etc.
* A topic has contents, that are written in Markdown.

This application implements two things: a web application and an API layer.

### The Web Application

The web application has a quite minimalistic design and it only consists of two
different pages:

1. The `login` page.
2. The `topics` page. This is the main page and in it the user can create new
   topics, read their contents, update them, etc.

The style for this web application has been taken from Reddit's one (the list
of topics in this application has the same style as the list of multireddits).
The markdown being shown has a similar style as the one from Github. A
snapshot:

![The main page](public/images/snapshot.png)

### The JSON API

This application also implements a JSON API. This API can be accessed by
setting `application/json` in the `Content-Type` or the `Accept` header of the
request. First of all, you have to login with the application. In order to do
so, you have to perform a POST HTTP request with the following body:

```json
{
  "name":     "name-of-the-user",
  "password": "password-for-this-user"
}
```

If everything was ok, the response should only be a JSON object with the
`token` key. The value for this key has to be used on every request afterwards
by adding it to the URL query part. With this in mind, we can now call any
method of this simple REST API:

|    Method   |        Path       | Keys in the JSON body |     Response    |
|:-----------:|:-----------------:|:---------------------:|:---------------:|
|GET          | /topics           |            -          | Array of Topics |
|POST         | /topics           |          name         |      Topic      |
|GET          | /topics/{topicId} |            -          |      Topic      |
|PUT or PATCH | /topics/{topicId} |    name or contents   |      Topic      |
|DELETE       | /topics/{topicId} |            -          |     Message     |


The fourth method is the update method. It accepts either the name or the
contents, but not both parameters at the same time. We might want to pass the
`name` key when renaming a topic. We will pass the `contents` key if we
want to update the contents of a topic. Moreover, the last method is the delete
method. This method responds with a `Message` object. A `Message` object
contains the key `msg` on success, and the `error` key on error.

Let's see a quick example (cURL with the `--trace-ascii` option):

    0000: POST /topics?token=6b6c0542-0891-4614-5dd8-92ce443dbcaf HTTP/1.
    0040: 1
    0043: User-Agent: curl/7.38.0
    005c: Host: localhost:3000
    0072: Accept: */*
    007f: Content-Type: application/json
    009f: Content-Length: 17
    00b3:
    => Send data, 17 bytes (0x11)
    0000: {"name": "topic"}

The create method will respond with the newly created Topic on success:

```json
{
  "id":         "ec654b88-e227-47bb-6254-60b77329732e",
  "name":       "topic",
  "contents":   "",
  "created_at": "2014-10-07T08:37:05.424276522+02:00",
  "markdown":   ""
}
```

So, let's explain the `Topic` object. The `id`, `name` and `created_at` columns
are quite self-explanatory. The `contents` column exposes the raw data of this
object. The `markdown` object contains the HTML code that has been produced
after rendering the markdown of the `contents` column.

## Getting this application up and running

This application uses the Go programming language and some awesome packages
like negroni, gorilla/mux, blackfriday, etc. Moreover, it uses PostgreSQL
for the database. In order to configure the database, you can use the following
environment variables:

- `TODO_DB_USER`, defaults to "postgres".
- `TODO_DB_NAME`, defaults to "todo-dev".
- `TODO_DB_PASSWORD`, defaults to "".
- `TODO_DB_HOST`, defaults to "localhost".
- `TODO_DB_SSLMODE`, defaults to "disable".

The port in which this application is listening into is set by the environment
variable `TODO_PORT`. By default it runs on the port 3000.

You might want to use the given **docker compose** setup as given in the
`docker-compose.yml` file. Note that it will create two nodes: `web` and `db`.
The `db` node uses the [official Postgres](https://hub.docker.com/_/postgres/)
image. Make sure to understand how it works. The `db` node has to be
initialized with the given command inside of the container:

    $ psql -U postgres < /tmp/db/tables.sql

There is also a `Dockerfile` providing an up-to-date image of the application
(`mssola/todo:latest` on the Docker Hub). However, if you want to build this
manually, you can type:

    $ godep go build && ./todo

### Secure connection

Since you usually want to run this through a safe connection, this application
also allows the deployer to set the following environment variables:

- `TODO_KEY_PATH`: the path to your key file.
- `TODO_CERT_PATH`: the path to the certificate file.

Moreover, the `.gitignore` file ignores the `docker-compose.production.yml`
file. You might want to use this file to run a customized version of
`docker-compose.yml` file. An example could be:

```yml
web:
  image: mssola/todo:latest
  volumes:
    - .:/go/src/github.com/mssola/todo
    - /path/to/certs:/path/to/certs
  ports:
    - 443:3000
  environment:
    TODO_DB_NAME: todo-production
    TODO_DB_HOST: todo_db_1
    TODO_KEY_PATH: /path/to/certs/todo.key
    TODO_CERT_PATH: /path/to/certs/todo.crt
  links:
    - db

db:
  image: library/postgres:9.4
  volumes:
    - ./db:/tmp/db
  environment:
    POSTGRES_DB: todo-production
```

## License

Copyright &copy; 2014-2015 Miquel Sabaté Solà

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

