# Code's Kitchen Backend

Code's Kitchen backend is a Golang backend using Golang 1.19

## Architecture

Following the Clean Architecture pattern and Layered Architecture

- service
- handler
- route
- cmd

## Service

Service contains business logic of how the appliances function

## Handler

Handler contains the HTTP handlers for calling the service layer

## Route

Route groups related handlers under the same endpoint

## Cmd

Cmd is serving as the dependency injection site and calls the main function calling the app

## Scripts

 `make install` for installing dependencies

 `make server` for running the application
