# Sock Shop
![go build](https://github.com/oshankkumar/sockshop/actions/workflows/go.yml/badge.svg)

A Golang demo web application. It is intended to aid the demonstration of how to write a web application in Golang using SOLID principles.
This is a re-write of https://github.dev/microservices-demo/ to demonstrate how we can use write a Go web app following SOLID principles.

## Overview

This is basic layout for Go web app.

## Package Layout

### `/api`
`api` package should only contain the request and response schema. It should not have any application logic or dependency. Additionally It can also contain interface definitions. Which can use the defined request and respose schema as input/output parmas in the interface methods.
It should not contain any concrete type which implements those interface. This will provide an abstraction to your business logic.

See [/api](./api/) package for example.

### `/cmd`

This will contain the main package(s) for this project. The directory name for each application in `cmd` should match the name of executable you want. You should not put a lot of code in main package. You should only initialize the app dependencies in here and inject it to higher level modules/components.

See [/cmd](./cmd/) package for example.

### `/domain`
`domain` package should only contain your application domain related information. It should define the application domain types and related interfaces. You should not put any concreate implementaion in here. This is an abstraction to your application domain. 
The actual implementation could be catered using a database or a third party API. Application domain should be independent of the actual concrete implementation.

See [/domain](./domain/) package for example.

### `/internal`

This package should provide the private application code. It should define the concrete types having methods with actual business logic.
Your application code can go in the `/internal/app` packgae.
It will contain concrete types which can implement intefaces defined in `/api` package.

You can define packages like `internal/db`, `internal/<api-client> ` ..etc. Which will define concrete types implementing interfaces defined in `/domain` package. 

See [/internal](./internal/) package for example. 

### `/transport`

This package should contain the application transport layer (eg ..http,grpc ..etc). You can define your application http Handlers in `/transport/http` package. 

See [/transport/http](./transport/http) package for example. 

