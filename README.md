# apibot

> Note: The current version is for preview only

[![](https://img.shields.io/badge/editor-code-2ca5e0?style=flat&logo=github)](https://github.com/pojol/apibot-editor)


# Try it out
Try the editor out [on website](http://1.117.168.37:7777/)

# Install
```shell
# run drive
$ docker pull braidgo/apibot:latest
$ docker run --rm -d  -p 8888:8888/tcp braidgo/apibot:latest
```

## Preview
[![image.png](https://i.postimg.cc/wT5HhYD3/image.png)](https://postimg.cc/6yQDXSjN)


### Control
* **Sequence** execute all child nodes under this node in sequence
* **Selector** Execute all nodes under this node, and exit this node when a child node is successfully executed

### Condition
* **Condition** Use expressions to determine whether to execute downward
* **Assert** Use expressions to determine whether to break execute

### Action
* **Http** Call an http request

### Decorator
* **Loop** Set the number of cycles of all child nodes under this node
* **Wait** Set a certain amount of time to wait at the current node

### API
* `/upload.blob`
* `/upload.file`
* `/get.list`
* `/get.blob`
* `/get.report`
* `/bot.create`
* `/bot.run`
* `/bot.step`
