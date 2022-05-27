# Golang workout-webservice

Simple Golang workout-webservice using:
* [Firebase](firebase.google.com) for authentication
* [Gin](github.com/gin-gonic/gin) for routing
* [Gorm](gorm.io/gorm) for ORM
* [Viper](github.com/spf13/viper) for environment settings

## Project Setup

You can rename the project by replacing all "workout-webservice" to your new project name in every file.

To start the project, you need to make sure that you have Go (Golang) and Docker installed. Then you can run these commands to install Golang dependency:

```
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get github.com/spf13/viper
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

### Integrating Firebase

If you want to use firebase, you need to duplicate `firebaseServiceAccountKey.json.sample` and remove the `.sample` to `firebaseServiceAccountKey.json` and provide your service account key in that file. Then, you need to run this command:

```
go get firebase.google.com/go
```

### Modify Configuration

This app will take configurations from environment variables or `.env` file in `/config`. If you're running this app locally, it is recommended to duplicate the `.env.sample` file in `/config` and rename it to `.env`.

### Integrating Database

You can change the database settings by editing the environment variables or `.env` file from the previous step. By default, this app will use PostgreSQL as database. If you want to run a PostgreSQL database in Docker, you can run:

```
bash start_db.sh
```

## Running the Program

You can run the program by simply running these commands:

```
go mod tidy
go run main.go
```

By default, you can access the app from `http://localhost:8080/`

## Future Plan

1. Integrate Redis
2. Websocket
3. Deploy to Google Cloud Platform example
4. Deploy to Amazon Web Service example
5. Deploy to Heroku example
6. Move Firebase integration to another branch
7. Add authentication provider in another branch