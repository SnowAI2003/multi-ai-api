FROM golang:1.22
WORKDIR /app

# Copy everything and build
COPY . .
RUN go mod tidy && go build -o server server.go

# Render expects this
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
