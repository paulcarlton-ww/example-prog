
FROM golang:1.13 as builder

WORKDIR /workspace

# copy modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache modules
RUN go mod download

# copy source code
COPY pkg/ pkg/
COPY api/ api/
COPY controllers/ controllers/
COPY internal/ internal/

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o example-prog pkg/main/main.go

FROM alpine:3.12

COPY --from=builder /workspace/example-prog /bin/

RUN apk -q add --no-cache --virtual .build-deps ca-certificates upx && \
    addgroup -S example-prog && \
    adduser -S -G example-prog example-prog && \
    upx -qqq /bin/example-prog && \
    apk -q del .build-deps

# Lock down system to example-prog user (no code changes).
USER example-prog

ENTRYPOINT ["example-prog"]
