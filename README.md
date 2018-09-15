# Chisk - Golang stdlib REST starter kit

[![Build Status](https://travis-ci.org/ribice/chisk.svg?branch=master)](https://travis-ci.org/ribice/chisk)
[![codecov](https://codecov.io/gh/ribice/chisk/branch/master/graph/badge.svg)](https://codecov.io/gh/ribice/chisk)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribice/chisk)](https://goreportcard.com/report/github.com/ribice/chisk)

Chisk - (go-chi starter kit), is a starter kit for building RESTful APIs using Golang stdlib only (with 3rd party libs being used for integrations, e.g. Redis, PostgreSQL).

It features only structured package management, middleware and reusable packages - [READ WHY](https://www.ribice.ba/beyond-rest).

Chisk's project structure mostly follows [THIS](https://github.com/golang-standards/project-layout) example repository.

## Wiring up new services

1. Create a new directory inside `/internal/` named after your domain, e.g. team. Inside it, create a file named `team.go` that will hold your main application service code (CRUD). Other application service code should be located in the same directory/package, named after its usage.

2. Directory `/internal/team/platform/` contains platform integrations, e.g. redis, postgresql, email, sms etc.

3. Directory transport contains http handlers, binders, request/response models, swagger annotations and other transport related code.

4. Wire up new service inside cmd/`service_name`/main.go, or use a single service ('monolith') like `api`.

## License

chisk is licensed under the MIT license. Check the [LICENSE](LICENSE.md) file for details.

## Author

[Emir Ribic](https://ribice.ba)