FROM golang:latest

WORKDIR /build

ENV GOPROXY https://goproxy.cn
COPY . .
RUN go mod download
RUN cd cmd && \
    GOARCH=amd64 GOOS=linux go build -o /output/backend .

FROM ubuntu:latest
WORKDIR /build
COPY --from=0 /output/backend /build/
RUN apt-get update \
    && apt-get install -y ca-certificates
EXPOSE 80
ENTRYPOINT ["/build/backend"]