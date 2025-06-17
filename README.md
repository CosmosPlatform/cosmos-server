# cosmos-server

Server for the cosmos platform

## File Structure

```
cosmos-server
├── pkg
│   ├── app                  // Package where the app is defined and has the functions to be initialized.
│   ├── api                  // Package where the API is defined
│   │   ├── routes           // Package where the routes are defined.
│   │   └── dto              // Package where the data transfer objects are defined.
│   ├── config               // Package where the configuration format is defined.
│   ├── server               // Package where the server is defined.
│   ├── services             // Package where the services are defined.
│   └── test                 // Package where test utilities are defined
└── config                   // Configuration files
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