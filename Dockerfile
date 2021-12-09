FROM golang:1.17.3-alpine3.14

## We copy everything in the root directory
## into a new workdir.
COPY . /app
WORKDIR /app

## we run go build to compile the binary
## executable of our Go program
RUN go build -v -o palindromee .

## Start binary.
ENTRYPOINT ["/app/docker-entrypoint.sh"]