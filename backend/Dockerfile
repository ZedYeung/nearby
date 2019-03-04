FROM golang:1.12.0-alpine3.9 as builder

WORKDIR /go/src
COPY ./src .

RUN apk add --no-cache git mercurial \
    && go get github.com/Masterminds/glide

RUN cd github.com/ZedYeung/nearby && glide install
RUN go build github.com/ZedYeung/nearby/main

# use a minimal alpine image
FROM alpine:3.9

# set working directory
WORKDIR /root
# copy the binary from builder
COPY --from=builder /go/src/main .
# run the binary
CMD ["./main"]