# syntax=docker/dockerfile:1

FROM golang:1.23 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/mstore

COPY go.mod ./
COPY go.sum ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /mstore

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /mstore /mstore

EXPOSE 8080
EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/mstore"]