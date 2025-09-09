FROM golang:1.23-alpine AS builder

ARG GITHUB_USER
ARG GITHUB_TOKEN
ARG GITHUB_PRIVATE_PATH

RUN apk add --no-cache git
RUN git config --global \
    url."https://$GITHUB_USER:$GITHUB_TOKEN@github.com/$GITHUB_PRIVATE_PATH".insteadOf \
    "https://github.com/$GITHUB_PRIVATE_PATH"

COPY go.mod go.sum Makefile /src/
COPY cmd/ /src/cmd/
COPY internal/ /src/internal/
COPY tools/ /src/tools/


WORKDIR /src
RUN go build -o /build/s3kd-service /src/cmd/s3kd
FROM harbor.infra.cluster.ionos.com/gcr.io/distroless/static:nonroot
COPY --from=builder /build/s3kd-service /bin/s3kd-service
USER nonroot
ENTRYPOINT ["/bin/s3kd-service"]
