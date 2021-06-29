# Developer Guide

## Getting Started

### Dependencies

Setup requires [Docker](https://docs.docker.com/install/).

Create a `.env` file in the root of this repo with following configurations:
```
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_HOST=
POSTGRES_PORT=4201
SERVER_IP=0.0.0.0
SERVER_PORT=80
ADMIN_USER=
ADMIN_PASSWORD=
SECRET_KEY=
```

`POSTGRES_HOST` should be the IP address of your computer in your connected network.

### Setting up the app

Checkout the branch for a given tutorial, and run `docker-compose build` or `docker-compose up --build`.

### Running the app on local machine

Run `docker-compose up` to see messages in the terminal. Run `docker-compose start` to run the app in the background. The app is available on http://localhost:4200.