FROM golang:1.19 AS build-stage

WORKDIR /app

ENV GOPROXY=https://goproxy.io,direct

COPY . .

COPY ./configs/envs/.env.prod ./.env 

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -v -o /main ./cmd/api/main.go

FROM gcr.io/distroless/base-debian11:latest-amd64 AS build-release-stage

WORKDIR /

COPY --from=build-stage /main /main
COPY --from=build-stage /app/.env /configs/envs/.env

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT [ "/main" , "serve-all"]