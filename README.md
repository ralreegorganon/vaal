# VAAL

The once and future skynet

## Use it

Create a new database user, and give it the password of ````ourrobotoverlords```` when prompted

    createuser -P vaal

Create a new database
    
    createdb -O vaal vaal

Set up the connection string
    
    export VAAL_CONNECTION_STRING="user=vaal password=ourrobotoverlords dbname=vaal sslmode=disable"

Install goose to run the database migrations
    
    go get bitbucket.org/liamstask/goose/cmd/goose

Run the migrations
    
    goose up

Run it locally
    
    go build
    ./vaal
