FROM golang:1.23.3-alpine AS builder

ARG APP_RELATIVE_PATH=./cmd/server/

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /app

COPY . .

# 设置镜像源
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy && go build -ldflags="-s -w" -o server ${APP_RELATIVE_PATH}

# 编译Go应用程序
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -a -installsuffix cgo -o server ${APP_RELATIVE_PATH}
# RUN GOOS=linux go build -ldflags="-s -w" -trimpath -a -installsuffix cgo -o ./server ${APP_RELATIVE_PATH}
# RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -trimpath -a -o server ${APP_RELATIVE_PATH}


FROM alpine:latest

# 设置工作目录
WORKDIR /app

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

ARG APP_CONF=config/prod.yaml
ENV OAI_APP_CONF=${APP_CONF}

COPY config ./config

COPY --from=builder /app/server .

# 设置可执行权限
RUN chmod +x ./server

# 设置环境变量（如果需要）
ENV TZ=Asia/Shanghai
ENV GIN_MODE=release

EXPOSE 8080
ENTRYPOINT [ "./server" ]

#docker build -t  1.1.1.1:5000/demo-api:v1 --build-arg APP_CONF=config/prod.yml --build-arg  APP_RELATIVE_PATH=./cmd/server/...  .
#docker run -it --rm --entrypoint=ash 1.1.1.1:5000/demo-api:v1
