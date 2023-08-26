FROM golang:1.21 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY src/ /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /ddns
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /app
COPY --from=build-stage /ddns /app
ENTRYPOINT ["/app/ddns"]