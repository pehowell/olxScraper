FROM golang:latest  as builder
RUN mkdir /build
ADD *.go /build/ 
WORKDIR /build 
RUN go get -d -v github.com/anaskhan96/soup
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o main .

FROM pehowell/alpine-dumbinit
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/main /app/
ENTRYPOINT ["/usr/local/bin/dumb-init", "/app/main"]
