FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git curl

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep

RUN mkdir -p $GOPATH/src/github.com/alexandrevilain/postgrest-auth
WORKDIR $GOPATH/src/github.com/alexandrevilain/postgrest-auth

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . ./
RUN go build -o postgrest-auth cmd/postgrest-auth/main.go

FROM alpine
WORKDIR /root
COPY --from=builder /go/src/github.com/alexandrevilain/postgrest-auth/postgrest-auth .
ENTRYPOINT ["./postgrest-auth"]