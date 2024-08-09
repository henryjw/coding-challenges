This is my solution for the [Build Your Own Memcached Server](https://codingchallenges.fyi/challenges/challenge-memcached/) challenge

## Usage
- Run `go run .` on this directory to run the code without saving an executable
  - Alternatively, you can run `go build .` to compile the code and generate an executable and then run `./memcached-server -p <port_number>`
    - Default port number is `9999`
    - Run `./memcached-server --help` to get more details on the supported parameters
- Follow link above for more details on this challenge

## Running tests
- Run `go test ./...` from this directory

## Features
- `get`, `set`, `add`, `delete`,  `replace`, `append`, and `prepend` commands
- Active deletion for expired cache entries
  - With this approach, expired data is periodically cleared
    - The frequency at which the background job runs is configurable in the code but is set to 1 second in the current implementation
- Passive deletion for expired cache entries
  - With this approach, expired data is only deleted when accessed