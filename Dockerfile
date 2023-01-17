FROM golang:1.19.3-alpine3.16

WORKDIR /challenge

COPY . .

RUN go build cmd/main.go

FROM alpine:3.14.0

WORKDIR /challenge

COPY --from=0 /challenge/main .
COPY --from=0 /challenge/internal/infra/persistence/postgresql/migrations ./migrations
EXPOSE 80
ENTRYPOINT [ "./main" ]



