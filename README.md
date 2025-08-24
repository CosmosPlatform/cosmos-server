# cosmos-server

Server for the cosmos platform

## File Structure

```
cosmos-server
├── config              // Local configuration files
├── api                   // Structs and validations of the server api requests and responses
├── db                   // Utilities and migrations for the postgress database
├── pkg
│   ├── app            // Package where the app is defined and has the functions to be initialized.
│   ├── config        // Package where there are structs and utilities for parsing configuration..
│   ├── errors        // Package where there are definitions for various kinds of general errors.
│   ├── log             // Package where the logger is defined.
│   ├── model        // Package where the various model of our application are defined..
│   ├── routes                  // Package where the route handlers for aour application are defined
│   ├── server               // Package where the server and its middleware is defined.
│   ├── services             // Package where the services are defined.
│   ├── storage             // Package where the storage service is defined.
│   └── test                 // Package where test utilities are defined
└── Makefile
```

## Requisites

### Golang

This program uses ``go 1.23.1`, so you should have go downloaded to run it.

### Migrations

This program needs a Mongo compatible database to work.

We use `go-migrate` so that we can easily keep the database in a valid state, and it can evolve side by side with this project's needs.

To download `go-migrate` run

```sh
sudo curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/migrate.list
sudo apt-get update
sudo apt-get install -y migrate
```

To run the migrations, first start a local database with `make run-db` and then run `make migrate`

### Environment Variables

This program expects a set of environment variables to be inserted in order to function properly.

- `JWT_SECRET`: Secret used to sign the JWT tokens.
  -  If you are running this locally you can quickly generate one running `openssl rand -base64 32`