# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage

WORKDIR /api

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /book-crud-api ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /book-crud-api /book-crud-api
COPY ./.env ./
COPY ./migrations/ ./migrations/

EXPOSE 7777

USER nonroot:nonroot

ENTRYPOINT [ "/book-crud-api" ]
