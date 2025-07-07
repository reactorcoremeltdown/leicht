FROM golang:alpine AS builder
RUN apk add git make build-base
COPY src/leicht /srv/leicht
WORKDIR /srv/leicht
ENV GOBIN=/usr/local/bin
ENV CGO_ENABLED=1
RUN go get && go build -o leicht

FROM alpine:3.21.3
COPY --from=builder /srv/leicht /srv/leicht
CMD [ "/srv/leicht/leicht" ]
