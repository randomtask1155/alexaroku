FROM golang:1.12 AS builder
LABEL maintainer="Daniel Lynch <danplynch@gmail.com>"
RUN mkdir -p /go/src/github.com/randomtask1155/alexaroku
WORKDIR $GOPATH/src/github.com/randomtask1155/alexaroku
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

ENV PORT=8080
ADD . .

RUN GOOS=linux GOARCH=arm GOARM=7 go build -o alexaroku .

FROM alpine:latest as certs
RUN apk --update add ca-certificates


FROM scratch
COPY --from=builder /go/src/github.com/randomtask1155/alexaroku/alexaroku /go/bin/alexaroku
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
#CMD ["./alexaroku"]
ENTRYPOINT ["/go/bin/alexaroku"]
