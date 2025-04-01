FROM golang:1.23.1-alpine AS build

WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o skinatar

FROM alpine:latest

COPY --from=build /app/skinatar /usr/local/bin/skinatar
RUN ls -l /usr/local/bin/skinatar
RUN chmod +x /usr/local/bin/skinatar

EXPOSE 8080
ENV REDIS_URL="localhost:6379"
VOLUME /app/cached_skins

CMD ["/usr/local/bin/skinatar"]