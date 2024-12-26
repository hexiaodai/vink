FROM --platform=$BUILDPLATFORM golang:1.22.7 AS build
ENV GOOS=linux \
    CGO_ENABLED=0 \
    GO111MODULE=on
WORKDIR /vink
COPY . .
RUN go build -o vink ./cmd/vink

FROM --platform=$BUILDPLATFORM alpine:3.16
COPY --from=build /vink/vink /bin/vink
RUN chmod go+x /bin/vink
ENTRYPOINT ["vink"]
