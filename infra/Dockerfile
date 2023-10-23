FROM golang:1.19 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o /gym-app-api


FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /gym-app-api /gym-app-api

EXPOSE 8080

USER nonroot:nonroot

CMD ["/gym-app-api"]