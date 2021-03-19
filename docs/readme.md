# Overview
[ACC Race Engineer server](https://github.com/blacksfk/are_server) (are_server) is an HTTP and UDP application that forwards data received from the [ACC Race Engineer client](https://github.com/blacksfk/acc_race_engineer) application to connected websocket clients observing the data via ACC Race Engineer interface (TODO: link repo).

## Project layout
This application takes suggestions from [here](https://github.com/golang-standards/project-layout#standard-go-project-layout) and [here](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) and uses them to implement the Clean architecture pattern with some minor changes for convenience due to the (small-ish) size of the project.

* `/` Types implementing the business logic of the application.

* `/cmd/are_server` Main application. Initialises services, creates routes, and injects dependencies.

* `/docs` Project documents describing the application.

* `/http` Controllers with methods implementing microframwork.Handler.

* `/http/middleware/validate` Validation logic in the form of middleware implementing microframwork.Middleware.

* `/mock` Mock types that implement interfaces defined in the business logic. Intended to be used for unit testing purposes.

## Business logic
TODO

## Database
TODO

## Sockets
TODO

## Handlers
TODO

## Routes
TODO
