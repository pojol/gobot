FROM alpine

COPY ./bot_linux /home/bot/
COPY ./script/* /home/bot/script/

WORKDIR /home/bot

# RUN go env -w GOPROXY=https://goproxy.cn

# RUN go build -buildmode=plugin /home/bot/plugins/json/json.go
# RUN go build -o bot_linux /home/bot/main.go

EXPOSE 8888

ENTRYPOINT ./bot_linux $0 $@