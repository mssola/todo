# ToDo [![Build Status](https://travis-ci.org/mssola/todo.svg?branch=master)](https://travis-ci.org/mssola/todo)

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

TODO

## The JSON API

TODO

## Getting this application up and running

This application uses the Go programming language and some awesome packages
like negroni, gorilla/mux, blackfriday, etc. Moreover, it uses PostgreSQL
for the database. You can install the dependencies manually and setup
the db/database.json file by yourself, or you can just perform the
following command:

    ./script/kickstart

It will ask for a few DB-related things: the user, the dbname and the password.
After this, you can install this like any other Go program and run it:

    go install
    todo

Last but not least, this application can be deployed to Heroku. Here are some
notes about it:

* This application uses PostgreSQL, so in you have to setup PostgreSQL in
  Heroku as well.
* Even though this application uses Sass, it has the resulting compiled CSS
  files in it. This is done this way so we don't have to mess with Sass on
  deployment.
* We'll have to use a third-party buildpack in order to deploy this
  application.

So, pick up your terminal and perform the following command:

    $ heroku create -b https://github.com/kr/heroku-buildpack-go.git

This will create the application with a Go buildpack setup. You should now go
to this application's config page on Heroku and add it the PostgreSQL add-on.
After this, no more configuration is needed, so you can just perform the
following command:

    $ git push heroku master

And that's it, you've deployed this application to Heroku!

## License

Copyright &copy; 2014 Miquel Sabaté Solà

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

