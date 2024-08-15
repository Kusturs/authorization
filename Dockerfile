FROM golang:1.22-alpine

# todo: refactor & check
# todo: test private dependencies here

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth_go /scr/cmd/main.go

EXPOSE 8080

CMD ["/auth_go"]