# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22.5
ARG ALPINE_VERSION=3.20
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH
ARG TARGETOS

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=$TARGETOS go build -ldflags="-w -s" -o /bin/server ./cmd

FROM gcr.io/distroless/static-debian12:nonroot AS final

COPY --from=build /bin/server /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]
