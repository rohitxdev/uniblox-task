# Golang API starter

<img src="./assets/public/images/go-fast.png" alt="Gopher flash" height="128" width="128"/>

## Requirements

- [Go 1.23.2+](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- Linux or macOS

## Setup

1. Clone the repository:

```bash
git clone https://github.com/rohitxdev/go-api-starter.git
```

2. Create a .env file in the root directory:

```bash
ENV=development
HOST=0.0.0.0
PORT=8080
SECRETS_FILE=secrets.json
```

3. Create a `secrets.json` file in the root directory:

```json
{
    "awsAccessKeyId": "",
    "awsAccessKeySecret": "",
    "googleClientId": "",
    "googleClientSecret": "",
    "sessionSecret": "",
    "databaseUrl": "",
    "smtpHost": "",
    "smtpUsername": "",
    "smtpPassword": "",
    "smtpPort": 465,
    "senderEmail": "",
    "s3Endpoint": "",
    "s3BucketName": "",
    "s3DefaultRegion": "",
    "allowedOrigins": [],
    "shutdownTimeout": "",
    "sessionDuration": "",
    "logInTokenExpiresIn": "",
    "jwtSecret": ""
}
```

## Development

Run the development server:

```bash
./run watch # or ./run docker-watch
```

## Production

Build the project:

```bash
./run build
```

Run the production server:

```bash
./run start
```

## Testing

Run the tests:

```bash
./run test
```

Generate a coverage report:

```bash
./run test-cover
```

Benchmark the project:

```bash
./run benchmark
```

## Notes

- The `run` script is used to automate common development/production tasks. Run `./run` to see the available tasks.
- The secrets should be stored in a `secrets.json` file in the root directory. If you want to use a different file, you can specify it using the `SECRETS_FILE` environment variable.
- HTTPS and Rate limiting should be handled at the reverse proxy level, not at the API level.
