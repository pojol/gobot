import React from "react";
import PubSub from "pubsub-js";
import Topic from "./topic";
import { append } from "@antv/x6/lib/util/dom/elem";
import { message } from "antd";
import OBJ2XML from "object-to-xml";
import Config from "./config";
import { Post, PostBlob } from "./request";
import Api from "./api";
import { NodeTy, IsScriptNode } from "./node_type";

export default class TreeModel extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      tree: {},
      botid: "",
      behaviorTreeName: "",
      httpCodeTmp: Config.httpCode,
      assertTmp: Config.assertCode,
      conditionTmp: Config.conditionCode,
    };
  }

  findNode = (id, parent, callback) => {
    if (parent.children && parent.children.length) {
      for (let nod of parent.children) {
        if (nod.id === id) {
          callback(parent, nod);
          break;
        }

        this.findNode(id, nod, callback);
      }
    }
  };

  addNode(parentid, childNode) {
    var syncinfo = (childNode) => {
      // 如果存在旧的节点参数数据，在这个行为里应该使用旧的数据；
      if (window.tree.get(childNode.id) === undefined) { 
        window.tree.set(childNode.id, childNode);
      }
    };

    console.info("add", childNode.ty, childNode.code)

    if (parentid === childNode.id) {
      // root
      this.setState({ tree: childNode }, () => {});
    } else {
      // init
      if (childNode.ty === NodeTy.Action && (childNode.code === "" || childNode.code === undefined)) {
        childNode.code = this.state.httpCodeTmp;
      } else if (childNode.ty === NodeTy.Condition && (childNode.code === "" || childNode.code === undefined)) {
        childNode.code = this.state.conditionTmp;
      } else if (childNode.ty === NodeTy.Assert && (childNode.code === "" || childNode.code === undefined)) {
        childNode.code = this.state.assertTmp;
      } else if (childNode.ty === NodeTy.Loop && childNode.loop === undefined) {
        childNode.loop = 1;
      } else if (childNode.ty === NodeTy.Wait && childNode.wait === undefined) {
        childNode.wait = 1;
      }

      if (this.state.tree.id === parentid) {
        var old = this.state.tree;
        old.children.push(childNode);
        this.setState({ tree: old }, () => {
          syncinfo(childNode);
        });
      } else {
        this.addChild(parentid, this.state.tree, childNode, (findChild) => {
          syncinfo(childNode);
        });
      }
    }
  }

  rmvNode(id) {
    this.findNode(id, this.state.tree, (parent, nod) => {
      parent.children.forEach(function (child, index, arr) {
        if (child.id === nod.id) {
          arr.splice(index, 1);
          if (window.tree.has(nod.id)) {
            window.tree.delete(nod.id);
          }
        }
      });
    });
  }

  rmvLink(id) {
    this.findNode(id, this.state.tree, (parent, nod) => {
      parent.children.forEach(function (child, index, arr) {
        if (child.id === nod.id) {
          arr.splice(index, 1);
        }
      });
    });
  }

  updateGraphInfo(graphinfo, notify) {
    this.findNode(graphinfo.id, this.state.tree, (parent, nod) => {
      nod.pos = graphinfo.pos;
    });
  }

  updateEditInfo(editinfo, notify) {
    window.tree.set(editinfo.id, editinfo);
    if (notify) {
      message.success("apply info succ");
    }
  }

  addChild = (findid, parent, child, callback) => {
    var flag = false;

    for (var i = 0; i < parent.children.length; i++) {
      if (parent.children[i].id == findid) {
        parent.children[i].children.push(child);
        flag = true;
        callback(child);
        break;
      }
    }

    if (!flag) {
      for (var i = 0; i < parent.children.length; i++) {
        this.addChild(findid, parent.children[i], child, callback);
      }
    }
  };

  syncNode = (parent, callback) => {
    for (var i = 0; i < parent.children.length; i++) {
      callback(parent.children[i]);

      if (parent.children[i].children && parent.children[i].children.length) {
        this.syncNode(parent.children[i], callback);
      }
    }
  };

  getTree() {
    var tree = this.state.tree;

    if (tree.children && tree.children.length) {
      this.syncNode(tree, (nod) => {
        var tar = window.tree.get(nod.id);

        if (window.tree.has(nod.id)) {
          if (IsScriptNode(nod.ty)) {
            nod.code = tar.code;
            nod.alias = tar.alias;
          } else if (nod.ty === NodeTy.Loop) {
            nod.loop = tar.loop;
          } else if (nod.ty === NodeTy.Wait) {
            nod.wait = tar.wait;
          }
        }
      });
    }

    return tree;
  }

  componentWillMount() {
    window.tree = new Map(); // 主要维护的是 editor 节点编辑后的数据
    this.setState({ tree: {} }); // 主要维护的是 graph 中节点的数据
    
    var remote = localStorage.remoteAddr
    if (remote === undefined || remote === "") {
      localStorage.remoteAddr = Config.driveAddr;
      remote = Config.driveAddr
    }
    window.remote = remote

    PubSub.subscribe(Topic.ConfigUpdate, (topic, info) => {
      if (info.key === "addr") {
        localStorage.remoteAddr = info.val;
        window.remote = info.val;
      } else if (info.key === "httpCode") {
        this.setState({ httpCodeTmp: info.val });
      } else if (info.key === "assertCode") {
        this.setState({ assertTmp: info.val });
      } else if (info.key === "conditionCode") {
        this.setState({ conditionTmp: info.val });
      }
      message.success("config update succ");
    });

    PubSub.subscribe(Topic.NodeAdd, (topic, linkinfo) => {
      this.addNode(linkinfo.parent, linkinfo.child);
    });

    PubSub.subscribe(Topic.NodeRmv, (topic, nodeid) => {
      this.rmvNode(nodeid);
    });

    PubSub.subscribe(Topic.LinkRmv, (topic, nodeid) => {
      this.rmvLink(nodeid);
    });

    PubSub.subscribe(Topic.UpdateNodeParm, (topic, info) => {
      this.updateEditInfo(info.parm, info.notify);
    });

    PubSub.subscribe(Topic.UpdateGraphParm, (topic, info) => {
      this.updateGraphInfo(info, false);
    });

    PubSub.subscribe(Topic.FileLoad, (topic, info) => {
      window.tree = new Map();
      this.setState({ tree: {}, behaviorTreeName: info.Name });
    });

    PubSub.subscribe(Topic.Create, (topic, info) => {
      var name = this.state.behaviorTreeName;
      var tree = this.getTree();

      if (name === undefined || name === "") {
        name = tree.id
        console.info("debug bot name", name)
      }

      var xmltree = {
        behavior: tree,
      };

      var blob = new Blob([OBJ2XML(xmltree)], {
        type: "application/json",
      });

      PostBlob(window.remote, Api.DebugCreate, name, blob).then(
        (json) => {
          if (json.Code !== 200) {
            message.error(
              "create fail:" + String(json.Code) + " msg: " + json.Msg
            );
          } else {
            this.setState({ botid: json.Body.BotID });
            message.success("create succ " + json.Body.BotID);
          }
        }
      )
    });

    PubSub.subscribe(Topic.Upload, (topic, filename) => {
      var tree = this.getTree();
      var xmltree = {
        behavior: tree,
      };

      var blob = new Blob([OBJ2XML(xmltree)], {
        type: "application/json",
      });

      PostBlob(window.remote, Api.FileBlobUpload, filename, blob).then(
        (json) => {
          if (json.Code !== 200) {
            message.error(
              "upload fail:" + String(json.Code) + " msg: " + json.Msg
            );
          } else {
            message.success("upload succ " + tree.id);
          }
        }
      );
    });

    PubSub.subscribe(Topic.Step, (topic, info) => {
      if (this.state.botid === "") {
        message.warn("have not created bot");
        return;
      }

      Post(window.remote, Api.DebugStep, { BotID: this.state.botid }).then(
        (json) => {
          if (json.Code !== 200) {
            message.warn(json.Msg);
            PubSub.publish(Topic.Blackboard, json.Body.Blackboard);
          } else {
            console.info("step", json.Body);
            PubSub.publish(Topic.Blackboard, json.Body.Blackboard);
            PubSub.publish(Topic.Focus, {
              Cur: json.Body.Cur,
              Prev: json.Body.Prev,
            });
          }
        }
      );
    });

    PubSub.subscribe(Topic.FileSave, (topic, msg) => {
      var tree = this.getTree();
      var xmltree = {
        behavior: tree,
      };

      var blob = new Blob([OBJ2XML(xmltree)], {
        type: "application/json",
      });

      // 创建一个blob的对象，把Json转化为字符串作为我们的值
      var url = window.URL.createObjectURL(blob);

      // 上面这个是创建一个blob的对象连链接，
      // 创建一个链接元素，是属于 a 标签的链接元素，所以括号里才是a，
      var link = document.createElement("a");

      link.href = url;

      // 把上面获得的blob的对象链接赋值给新创建的这个 a 链接
      // 设置下载的属性（所以使用的是download），这个是a 标签的一个属性
      link.setAttribute("download", "behaviorTree.xml");

      // 使用js点击这个链接
      link.click();
    });
  }

  componentDidMount() {}

  render() {
    return <div></div>;
  }
}
