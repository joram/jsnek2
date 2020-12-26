FROM golang:1.14.3-alpine as builder
RUN apk add git
WORKDIR /go/src/github.com/joram/jsnek2/
RUN go get github.com/davecgh/go-spew/spew
RUN go get github.com/julienschmidt/httprouter
RUN go get github.com/BattlesnakeOfficial/rules
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o jsnek2 .
RUN ls -hal

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /jsnek2/
COPY --from=builder /go/src/github.com/joram/jsnek2/jsnek2 .
RUN chmod +x /jsnek2/jsnek2
CMD ["/jsnek2/jsnek2"]