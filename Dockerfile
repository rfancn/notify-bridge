FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git
WORKDIR /src

COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w \
    -X github.com/rfancn/notify-bridge/cmd.version=${VERSION} \
    -X github.com/rfancn/notify-bridge/cmd.commit=${COMMIT} \
    -X github.com/rfancn/notify-bridge/cmd.buildDate=${BUILD_DATE}" \
    -o /out/notify-bridge .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /out/notify-bridge /usr/local/bin/notify-bridge
COPY config /app/config

EXPOSE 8080
ENTRYPOINT ["notify-bridge", "run"]