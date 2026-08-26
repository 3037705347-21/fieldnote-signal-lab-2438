FROM golang:1.26.5
WORKDIR /app
COPY . .
RUN go build ./...
EXPOSE 18081
CMD ["go", "run", "./cmd/api"]
