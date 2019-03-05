FROM golang:1.12.0-alpine3.9 as builder

WORKDIR /go
COPY . .

RUN apk add --no-cache git mercurial \
    && go get github.com/Masterminds/glide

RUN cd src/github.com/ZedYeung/nearby && glide install
RUN go build github.com/ZedYeung/nearby/main

# use a minimal alpine image
FROM alpine:3.9

# set working directory
WORKDIR /root
# copy the binary from builder
COPY --from=builder /go/main .
COPY ./GOOGLE_APPLICATION_CREDENTIALS.json .

# https://stackoverflow.com/questions/52341878/cannot-exchange-accesstoken-from-google-api-inside-docker-container
# to use https inside alpine
RUN apk add --no-cache ca-certificates openssl

# run the binary
CMD ["./main"]