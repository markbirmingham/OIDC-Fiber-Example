# OIDC Fiber Example

Example OIDC integration with the Fiber web framework (https://gofiber.io), adapted from the excellent Auth0 guide (https://auth0.com/docs/quickstart/webapp/golang), which is based on the _Gin Web Framework_. I have simplified the file structure somewhat.

This toy app uses the CoreOS OIDC library (https://github.com/coreos/go-oidc), and it's compatible with _Keycloak_.

# Instructions

Copy the `example.env` template to be the default (`.env`):

```
cp example.env .env
```

Update the values in the new `.env` file with OIDC configuration.

Download the _Go_ modules:

```
go mod download
```

Run the app:

```
go run .
```

Open a browser and go to http://localhost:9000

Build a static binary:

```
go build
```
