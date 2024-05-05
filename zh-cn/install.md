# 安装 gobot
  
---

## 本地模式
1. 进入最新的 [release](https://github.com/pojol/gobot/releases/tag/v0.4.4)页面 下载可执行程序
2. 执行 gobot_driver_win_x64_v0.4.4 目录中的 start.bat 文件， 运行服务器
3. 执行 gobot_editor_win_x64_v0.4.4 目录中的 gobot.exe， 运行编辑器程序
4. 在弹出的地址输入窗口 或 config 页的地址栏中填入 http://127.0.0.1:8888 本地服务器地址

---

## 服务器模式（使用 Docker-compose 安装到服务器
1. 安装 docker-compose
    ```shell
    # for CentOS
    yum install docker-compose -y

    # for Ubuntu
    apt-get install docker-compose -y
    ```

2. 下载并编辑 [docker-compose.yml](https://github.com/pojol/gobot-driver/blob/develop/docker-compose.yml) 在启动前，请确保对 `MYSQL_ROOT_PASSWORD` 和 `MYSQL_PASSWORD` 参数进行赋值
3. 运行命令 `docker-compose up -d` 运行成功后，访问 http://localhost:7777/ 即可进行 gobot 的编辑

---

## 服务器模式（使用 k8s 进行集群式部署