# build stage
FROM golang:latest AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./internal/cmd/app/main.go
LABEL authors="leoscrowi"

# run go app stage
FROM scratch
COPY --from=build /src/main /main
COPY --from=build /src/migrations /migrations
EXPOSE 8080
CMD ["/main"]