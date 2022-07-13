# mini-wallet
Service that handle the mini wallet

# Description
Mini wallet service has 4 folders, such as:
1. Controllers: Contains routeController.go that act as bridging between main and logics. No business logic happens here
2. workers: Contains files that handle all the business logic such as validation, update datbase, connect to database, view database
3. models: Contains files that is representative of database. The file name is the repersentative of data
4. database: Contains json files that act as the database. We use json file to view, update and delete data


# How to run
1. Install golang
2. Make new folders in GOPATH: pkg, src, bin
3. Clone this repo into src folder using command git clone https://github.com/bylardo/mini-wallet.git
4. Run command go mod init miniwallet.co.id
5. Run command go get github.com/gorilla/mux
6. Run command go build
7. Run command go run .
8. Open http://localhost:1991