FROM golang:1.21 AS builder
ARG GOPROXY
ENV GOPROXY $GOPROXY
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o naaprs cmd/naaprs/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/naaprs .

CMD ["./naaprs"]
