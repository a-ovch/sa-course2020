FROM golang:1.14 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:3.11.5
RUN adduser -D healthcheck
USER healthcheck
WORKDIR /app
COPY --from=builder /app/sa-course-app .
EXPOSE 8000
ENV SERVICE_PORT 8000
CMD ["./sa-course-app"]