FROM golang:1.14 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:3.11.5
RUN adduser -D app-executor
USER app-executor
WORKDIR /app
COPY --from=builder /app/sa-course-app .
CMD ["./sa-course-app"]