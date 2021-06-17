# Roulette

## Configure
The following elements are accepted as environment variables in `./settings.yaml`
```yaml
port: ":8080" // The port on which the server is to run.
```

## Build & Run
This will lint the codebase, create a binary and start the server.
```makefile
make execute
```

## Design
This project was designed with
- [Twelve-Factor](https://12factor.net/) in mind
- A commonly adopted project [structure](https://github.com/golang-standards/project-layout)
- Golang's official code review [guide](https://github.com/golang/go/wiki/CodeReviewComments)

