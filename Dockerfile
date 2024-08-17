FROM golang:1.22.5-alpine AS builder

WORKDIR /src

COPY ./go.* ./
RUN go mod download

COPY ./ ./

RUN GOARCH=arm64 CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

FROM scratch AS final

ENV WAIT_STARTUP_TIME 15
ENV WAIT_LIVENESS_TIME 20
ENV WAIT_READINESS_TIME 20

COPY --from=builder /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
