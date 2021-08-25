FROM alpine

COPY ./bot_linux /home/app/
COPY ./plugins/json/json.so /home/app

WORKDIR /home/app

EXPOSE 8888

ENTRYPOINT ./bot_linux $0 $@