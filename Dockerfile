FROM golang:alpine AS builder

ARG GOOS=linux
ARG GOARCH=amd64
ARG GOARM
ARG GOPROXY
ARG LDFLAGS

ENV GOPROXY $GOPROXY

WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build -ldflags "${LDFLAGS}" -o naaprs cmd/naaprs/main.go

FROM scratch
COPY --from=builder /app/naaprs /naaprs

CMD ["/naaprs"]
