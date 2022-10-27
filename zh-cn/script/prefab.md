# 将自己的预定义节点引入到bot中

* 通过文件
* 通过配置
* 通过节点

### 通过文件
> 通过文件的方式需要自己构建gobot的server端

```shell
# 1. 
git clone https://github.com/pojol/gobot

# 2. bot 在执行前会将 /script 目录下的 .lua 文件加载到执行环境中
# 将自己的脚本copy到script目录

# 3.
# 通过 docker build 构建自己的镜像
```

### 通过配置
> 在 editor/config/#Prefab script node# 栏中添加自己的预定义节点
* global 节点（默认
    - 这个系统提供的内置预定义节点会在每次执行bot时加载到执行环境中
* script 节点
    - 自行可添加的节点都默认是 script 节点，在 apply 之后节点都会在 editor/edit 面板的 prefab 栏中出现， 通过拖拽到页面的形式引用

### 通过节点（不推荐的
> 在root节点下添加脚本节点，并将自己的脚本填写在其中就可以啦；

> 通过这种方式可以免去自己构建工程，但是相应的在每颗树上面都要操作挂接一下自己的脚本；

![](../../res/preload.png)