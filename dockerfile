FROM golang:alpine as build

RUN mkdir /src && \
      apk update && \
      apk add git

WORKDIR /src
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/monex .

FROM alpine:3.18.0
RUN apk --no-cache add ca-certificates

ARG BUILD_DATE

LABEL \
    com.datadoghq.ad.check_names='["openmetrics"]' \
    com.datadoghq.ad.init_configs='[{}]' \
    com.datadoghq.ad.instances='[{"openmetrics_endpoint":"http://monex:9090/metrics","namespace":"example","metrics":["example*"]}]' \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.authors="<somebody>" \
    org.opencontainers.image.url="<somewhere>" \
    org.opencontainers.image.source="https://github.com/twistingmercury/monitoring-examples" \
    org.opencontainers.image.version="1.0.0" \
    org.opencontainers.image.title="TwistingMercury monitoring-packages example" \
    org.opencontainers.image.description="TwistingMercury monitoring-packages example"

EXPOSE 8080 9090

ENV OTEL_COLLECTOR_EP=""

WORKDIR /app
COPY --from=build /src/bin/ /app/

ENTRYPOINT [ "./monex"]