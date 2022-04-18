FROM golang:1.18-alpine AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# ARG VERSION=n/a \
#     BUILD_DATE=n/a

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -o ns-scheduler .


FROM alpine:3.15.4

WORKDIR /app

COPY --from=builder /build/ns-scheduler .

ENTRYPOINT [ "/app/ns-scheduler" ]