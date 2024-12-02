FROM --platform=$BUILDPLATFORM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /main .

# TODO: It can be probably FROM scratch if I copy ca-certs(?)
FROM golang:latest

WORKDIR /

COPY --from=builder /main .
EXPOSE 8080

CMD ["/main"]
