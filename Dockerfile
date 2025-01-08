FROM alpine:3.21.2 AS certs

RUN apk add ca-certificates

FROM node:lts as tob

WORKDIR /build
COPY . /build

RUN npm ci
RUN npm run build

FROM golang:1.22 AS gob

ARG VERSION

WORKDIR /build

COPY --from=tob /build .

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-X main.version=${VERSION}" -a -o prompage main.go

FROM alpine:3.21.2

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gob /build/prompage /prompage

ENTRYPOINT [ "/prompage" ]
