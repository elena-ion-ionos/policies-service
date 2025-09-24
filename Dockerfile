FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

RUN --mount=type=secret,id=github_token \
    git config --global \
      url."https://x-access-token:$(cat /run/secrets/github_token)@github.com/ionos-cloud".insteadOf \
      "https://github.com/ionos-cloud"

COPY go.mod go.sum Makefile /src/
COPY cmd/ /src/cmd/
COPY internal/ /src/internal/

WORKDIR /src
RUN go build -o /build/go-service /src/cmd/service
FROM harbor.infra.cluster.ionos.com/gcr.io/distroless/static:nonroot
COPY --from=builder /build/go-service /bin/go-service
USER nonroot
ENTRYPOINT ["/bin/go-service"]
