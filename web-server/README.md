This is my solution for the [Build Your Own Web Server](https://codingchallenges.fyi/challenges/challenge-webserver/) challenge

## Usage
- Run `go run .` on this directory to run the code without saving an executable
    - Alternatively, you can run `go build .` to compile the code and generate an executable and then run `./memcached-server -p <port_number>`
        - Default port number is `8080`
        - Run `./web-server --help` to get more details on the supported parameters
- Follow link above for more details on this challenge

## Running tests
- Run `go test ./...` from this directory

## Functionality
- Returns HTML from `www` folder
- If a file that doesn't exist is requested, the server will return a 404 error
- Handles concurrent requests