# deploy-builder Stage
FROM golang:1.19-alpine AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app cmd/todo/main.go


# deploy Stage
FROM alpine:latest AS deploy

RUN apk update

COPY --from=deploy-builder /app/app .

CMD ["./app"]


# dev Stage
FROM golang:1.19 as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

CMD ["air"]
