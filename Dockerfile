FROM alpine

COPY ./bot_linux /home/
COPY ./script/* /home/script/

WORKDIR /home

# RUN go env -w GOPROXY=https://goproxy.cn

RUN echo -e "https://mirrors.ustc.edu.cn/alpine/latest-stable/main\nhttps://mirrors.ustc.edu.cn/alpine/latest-stable/community" > /etc/apk/repositories && \
    apk update &&\
    apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" >  /etc/timezone

EXPOSE 8888
EXPOSE 6060

ENTRYPOINT ./bot_linux $0 $@