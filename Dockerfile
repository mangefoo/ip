############################
# STEP 1: Build executable
############################
FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

RUN adduser -D -g '' appuser

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app .

############################
# STEP 2: Build image
############################
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /go/bin/app /app

USER appuser

ENTRYPOINT ["/app"]
