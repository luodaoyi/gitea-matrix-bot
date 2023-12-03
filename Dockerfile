FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 1
#ENV GOPROXY https://goproxy.cn,direct
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -a -installsuffix cgo -o /app/gitea-matrix-bot .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/gitea-matrix-bot /app/gitea-matrix-bot
# /app 加入到path中
ENV PATH /app:$PATH
COPY config.ini.example /app/config.ini
CMD ["./entrypoint.sh"]