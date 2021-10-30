FROM nginx:1.21.3

COPY ./build/ /usr/share/nginx/html/
RUN ls -la /usr/share/nginx/html/*
COPY ./config/nginx_default.conf /etc/nginx/conf.d/default.conf

WORKDIR /home

EXPOSE 7777

CMD ["nginx","-g","daemon off;"]
