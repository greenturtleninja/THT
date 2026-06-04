# THT

A REST API for Eagle Bank which conforms to this OpenAPI specification which allows a user to create, fetch, update and delete a bank account and deposit or withdraw money from the account. These will be stored as transactions against a bank account which can be retrieved but not modified or deleted.

## Application

folder

\api<br />
- main.go -> entry point for the api<br />
<br />
\controller<br />
-- \Account\User....
    -- controller code responsible for fetching and returning data
\models<br />
-- \user\account\....<br />
    -- model code for queryingthe data source and retuning the data object
<br />
\auth<br />
-- \JWT<br />
    -- code for creating and checking the jwt token<br />

## Running the application

### Create a new token

Token keys with rsa using certs dir
```
> cd certs
> openssl genrsa -out priv_key.pem 2048
> openssl rsa -pubout -in priv_key.pem -out pub_key.pem

```

### run the application
```
go run api/main.go
```
This should start the server.


## Postgres

Install postgress and update the postgres database with the table schemas in database folder

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


