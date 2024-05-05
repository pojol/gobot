> Gobot是一个有状态的api/协议测试工具，支持图形化的行为编辑/调试、脚本节点、压力测试、测试报告等


### 特性
* 使用`行为树`控制机器人的运行逻辑，使用`脚本`控制节点的具体行为（比如发起一次http请求
* 提供图形化的编辑，调试能力
* 可以`预制`模版节点，在编辑器中直接使用预制过的节点（可通过标签筛选
* 可以通过 http api `'curl post /bot.run -d '{"Name":"某个机器人"}'` 驱动一个阻塞式的机器人，通过这种方式可以方便的集成进`CI`中的测试流程
* 可以进行`压力测试`（可以在配置页设置并发数
* 提供压力测试后的API/协议`报告`查看

> 注: 演示视频录制的版本较老，但主要的使用逻辑没变
【视频演示】 https://www.bilibili.com/video/BV1sS4y1z7Dg/?share_source=copy_web&vd_source=e3fcc3b3ff4b88affe4fd8529c90e1c2

### 预览
![img](/res/preview.png)