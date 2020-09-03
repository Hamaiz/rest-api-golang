# REST API with golang :)
This is a REST API that i made with golang. I used the cli tool called cobra for making the api run with cli commands. At this point it uses:
* Gorilla sessions for creating secure sessions
* MongoDB to store sessions
* Gorilla mux for routing
* Basic net/smtp package for validating the email address
* PostgreSQL for storing users information

I did alot of diging online in trying to make this REST API architecture design. I used cobra because i saw other people using the same style.

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. You have to have Golang installed on your machine and have basic understanding of the language.
### Installing
1. Clone the repo
```
git clone https://github.com/Hamaiz/.git
```
2. Get Cobra
```
go get -u github.com/spf13/cobra/cobra
```
3. Create .env files
```
touch .env.development
```
5. Add all the enviornment variables form the **.env.local** file to **.env.development**
* DATABASE_URL= (PostgreSQL database url)
* SESSION_DB= (MongoDB database url)
* SESSION_KEY= (Random session key)
* GM_EMAIL= (Gmail email)
* GM_PASS= (Gmail password)

6. Run main.go file
```
go run main.go
```
7. Start the api server
```
go run main.go serve
```
7. Explore
```
Enjoy! :)
```
## Built With
* [net/http](https://golang.org/pkg/net/http/) - Package provides basic http client and server implementation
* [net/smtp](https://golang.org/pkg/net/smtp/) - Package provides Simple Mail Transfer Protocol (smtp)
* [cobra](https://github.com/spf13/cobra) - Tool in golang for creating cli tools
* [gorilla/sessions](https://github.com/gorilla/sessions) - Session implementaion with golang
* [gorilla/mux](https://github.com/gorilla/mux) - Basic router mux
* [MongoDB](https://github.com/globalsign/mgo) - MongoDB driver for golang
* [PostgreSQL](https://github.com/jackc/pgx) - PostgreSQL driver for golang
## Contributor
* **Ali Hamaiz** - [Portfolio](https://thanksdear.herokuapp.com/)
