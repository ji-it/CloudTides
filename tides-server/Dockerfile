FROM golang:latest

ENV GOPROXY https://goproxy.cn
COPY . .
RUN go mod download
RUN mkdir /output && \
    cd cmd && \
    GOARCH=amd64 GOOS=linux go build -o /output/backend .

EXPOSE 80
ENTRYPOINT ["/output/backend"]