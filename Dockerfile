FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

COPY go.mod go.sum Makefile /src/
COPY cmd/ /src/cmd/
COPY internal/ /src/internal/

WORKDIR /src

RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
    git config --global \
      url."https://x-access-token:$(cat /run/secrets/GIT_AUTH_TOKEN)@github.com/ionos-cloud".insteadOf \
      "https://github.com/ionos-cloud"

RUN go build -o /build/go-service /src/cmd/service
FROM harbor.infra.cluster.ionos.com/gcr.io/distroless/static:nonroot
COPY --from=builder /build/go-service /bin/go-service
USER nonroot
ENTRYPOINT ["/bin/go-service"]
