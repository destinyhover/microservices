FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /payment ./cmd/

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /payment /app/payment
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/app/payment"]