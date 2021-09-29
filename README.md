# apibot

> Note: The current version is for preview only

[![](https://img.shields.io/badge/editor-code-2ca5e0?style=flat&logo=github)](https://github.com/pojol/gobot-editor)


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



### API
* `/file.txtUpload`
* `/file.blobUpload`
* `/file.remove`
* `/file.list`
* `/file.get`

* `/bot.create`
* `/bot.list`
* `/bot.info`

* `/debug.create`
* `/debug.step`


### Script
* `http.post`
* `http.get`
* `http.put`

* `json.encode`
* `json.decode`

* `proto.marshal`
* `proto.unmarshal`