FROM alpine:3.19.1 AS certs

RUN apk add ca-certificates

FROM golang:1.22 AS builder

WORKDIR /build

COPY . /build

RUN go mod download
RUN CGO_ENABLED=0 go build -a -o prompage main.go

FROM alpine:3.19.1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/prompage /prompage

ENTRYPOINT [ "/prompage" ]
