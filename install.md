# Install

---

## Local Mode
1. Go to the latest [release](https://github.com/pojol/gobot/releases/tag/v0.4.4) page to download the executable program.
2. Execute the start.bat file in the `gobot_driver_win_x64_v0.4.4` directory to run the server.
3. Execute gobot.exe in the `gobot_editor_win_x64_v0.4.4` directory to run the editor program.
4. In the popped-up address input window or the address bar on the config page, enter the local server address http://127.0.0.1:8888.

---

## Server Mode (Installing on Server using Docker-compose)
1. install docker-compose
    ```shell
    # for CentOS
    yum install docker-compose -y

    # for Ubuntu
    apt-get install docker-compose -y
    ```
2. Download and edit [docker-compose.yml](https://github.com/pojol/gobot-driver/blob/develop/docker-compose.yml). Before starting, make sure to assign values to the MYSQL_ROOT_PASSWORD and MYSQL_PASSWORD parameters.
Run the command docker-compose up -d. After successful execution, access http://localhost:7777/ to edit Gobot.

--- 

## Server Mode (Cluster Deployment using k8s)
