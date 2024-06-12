FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/authors/main.go
FROM golang
WORKDIR /app
COPY --from=builder --chown=${USERNAME}:${USERNAME} /app/main .
CMD [ "/app/main" ]