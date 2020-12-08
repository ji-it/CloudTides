FROM golang:latest

WORKDIR /build

ENV GOPROXY https://goproxy.cn
COPY . .
RUN go mod download
RUN cd cmd && \
    GOARCH=amd64 GOOS=linux go build -o /output/backend .

FROM ubuntu:latest
COPY --from=0 /output/backend .
EXPOSE 80
ENTRYPOINT ["./backend"]