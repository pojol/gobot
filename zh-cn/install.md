# 安装


### 本地无数据库模式安装
> 可以用于体验，无依赖下载即可使用

1. 进入 [release地址](https://github.com/pojol/gobot/releases/tag/v0.3.8) 下载相应资源
2. 执行 gobot_driver_win_x64_v0.3.8 目录中的 run.bat 文件， 运行服务器
3. 执行 gobot_editor_win_x64_v0.3.8 目录中的 gobot.exe， 运行编辑器程序
    * 在弹出的地址输入窗口 或 config 页的地址栏中填入 http://127.0.0.1:8888 本地服务器地址
4. 切换到编辑器的 bots 面板，将 `http_sample.txt` 和 `tcp_sample.txt` 两个用例拖入
5. 选择一个用例，点击 load 将机器人加载到编辑界面
    * 点击下方的 debug （爬虫）按钮进行调试（创建一个新的调试机器人
    * 点击旁边的 运行 按钮，单步执行（运行行为树节点
    * 点击编辑器中的任意一个节点 可以查看这个节点的设置
    * Meta 面板 可以查看机器人的所有数据
    * Response 显示的是每个节点中的返回值
    * RuntimeErr 显示的是执行节点可能遇到的错误信息（会自动切换过去


### Docker-compose 安装
> 可以通过网页访问，数据持久化在 mysql
1. 安装 docker-compose
    ```shell
    # for CentOS
    yum install docker-compose -y

    # for Ubuntu
    apt-get install docker-compose -y
    ```

2. 下载并编辑 [docker-compose.yml](https://github.com/pojol/gobot-driver/blob/develop/docker-compose.yml) 在启动前，请确保对 `MYSQL_ROOT_PASSWORD` 和 `MYSQL_PASSWORD` 参数进行赋值
3. 运行命令 `docker-compose up -d` 运行成功后，访问 http://localhost:7777/ 即可进行 gobot 的编辑


### K8s安装
> 可以通过网页访问，数据持久化在 mysql，可以开启分布式模式，添加任意个 drive 节点支撑巨量机器人测试需求
