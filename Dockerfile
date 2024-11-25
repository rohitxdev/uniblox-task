# syntax=docker/dockerfile:1
ARG BASE_IMAGE_TAG=:1.23.2-alpine

# Development image
FROM golang${BASE_IMAGE_TAG} AS development

WORKDIR /app

RUN apk add --no-cache build-base bash git && go install github.com/air-verse/air@latest

RUN go mod download -x

ENTRYPOINT ["./run","watch"]


# Production builder image
FROM golang${BASE_IMAGE_TAG} AS production-builder

WORKDIR /app

RUN apk add --no-cache build-base bash git

ENV GOPATH=/go

COPY go.sum go.mod ./

RUN go mod download -x && go mod verify

COPY . .

RUN ./run build


# Final production image
FROM scratch AS production

WORKDIR /app

COPY --from=production-builder /app/bin .

ENTRYPOINT ["./main"]