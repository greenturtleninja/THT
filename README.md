# THT

A REST API for Eagle Bank which conforms to this OpenAPI specification which allows a user to create, fetch, update and delete a bank account and deposit or withdraw money from the account. These will be stored as transactions against a bank account which can be retrieved but not modified or deleted.

## Application

folder

\api
- main.go -> entry point for the api

\models
-- \user
    -- struct for user CRUD operations
-- \account
    -- struct for account CRUD operations
-- \transaction
    -- struct for transaction CRUD operations

\pkg
-- \jwt
    -- code for creating and checking the jwt token

## Running the application

First run and install sqlite 3 - instructions below

```
go run api/main.go
```
This should start the server. Sqlite3 doesn't need to be running just needs to be setup


## SQLite3

Install sqlite3 as per instructions - https://sqlite.org/

From the root directory run the below from the command line

```
>sqlite3 eagle_bank.db
sqlite>.database
sqlite>.read database\schema\accounts.sql
... repeat for all the schemas
sqlite>.quit

Also the below is very helpful for listing specific sqlite commands
sqlite>.help 
```

## TODO

Move database connection to init function

Implement JWT AUTH and return user for accessing users accounts

Validation on requests

Finish of api calls with attached to hello function

Unit tests

Condense the handler functions

Clean up code, lots of repeats and maybe typos

Database checks for valid repeated ids

Switch to Mysql/or some other SQL


