# 将自己的脚本添加到行为树中

* 通过文件
* 通过节点

### 通过文件
> 通过文件的方式需要自己构建gobot的server端

```shell
# 1. 
git clone https://github.com/pojol/gobot

# 2. 
# 将自己的脚本copy到script目录

# 3.
# 通过 docker build 构建自己的镜像
```

### 通过节点
> 在root节点下添加脚本节点，并将自己的脚本填写在其中就可以啦；

> 通过这种方式可以免去自己构建工程，但是相应的在每颗树上面都要操作挂接一下自己的脚本；

![](../../res/preload.png)