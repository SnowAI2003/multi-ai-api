FROM golang:1.22
WORKDIR /app

COPY . .
RUN go mod tidy && go build -o server server.go

ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
