import React from "react";
import PubSub from "pubsub-js";
import Topic from "./topic";
import { message } from "antd";
import OBJ2XML from "object-to-xml";
import Config from "./config";
import { Post, PostBlob } from "./request";
import Api from "./api";
import { NodeTy, IsScriptNode } from "./node_type";


/*!

  // relation info
  {
    id : string,
    children : []
  }

  // node info
  {
    id : string // node id
    ty : string // node type
    pos : {
      x : number,
      y : number,
    },
    code : "",
    wait : 0,
    loop : 0,
    children : []
  }

*/

export default class TreeModel extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      rootid: "",
      nods: [], //  root 记录节点的链路关系， window(map 记录节点的细节
      botid: "",
      behaviorTreeName: "",
      httpCodeTmp: Config.httpCode,
      assertTmp: Config.assertCode,
      conditionTmp: Config.conditionCode,
    };
  }

  _getRelationInfo(parentChildren, children) {

    var cinfo = {
      id: children.id,
      children: []
    }

    parentChildren.push(cinfo)

    if (children.children && children.children.length) {
      children.children.forEach(cc => {
        this._getRelationInfo(cinfo.children, cc)
      })
    }

  }

  getRelationInfo(nod) {
    var rinfo = {
      id: nod.id,
      children: []
    }

    if (nod.chlidren && nod.children.length) {
      nod.children.forEach(children => {
        this._getRelationInfo(rinfo.children, children)
      });
    }

    return rinfo
  }

  setNode(nod) {
    // init
    if (nod.ty === NodeTy.Action && (nod.code === "" || nod.code === undefined)) {
      nod.code = this.state.httpCodeTmp;
    } else if (nod.ty === NodeTy.Condition && (nod.code === "" || nod.code === undefined)) {
      nod.code = this.state.conditionTmp;
    } else if (nod.ty === NodeTy.Assert && (nod.code === "" || nod.code === undefined)) {
      nod.code = this.state.assertTmp;
    } else if (nod.ty === NodeTy.Loop && nod.loop === undefined) {
      nod.loop = 1;
    } else if (nod.ty === NodeTy.Wait && nod.wait === undefined) {
      nod.wait = 1;
    }

    window.tree.set(nod.id, nod)
  }

  syncMapInfo(nod) {

    this.setNode(nod)

    if (nod.children && nod.children.length) {
      nod.children.forEach(children => {
        this.syncMapInfo(children)
      })
    }

  }

  addNode = (nod) => {

    if (nod.ty === NodeTy.Root) {
      this.setState({ rootid: nod.id })
    }

    let rinfo = this.getRelationInfo(nod)
    this.syncMapInfo(nod)

    let olst = this.state.nods
    olst.push(rinfo)

    this.setState({ nods: olst })
  }

  rmvNode = (id) => {

    if (id === this.state.rootid) {
      return
    }

    /*
    this.state.nods.forEach(children => {
      this.findNode(children.id, this.state.tree, (parent, nod) => {
        parent.children.forEach(function (child, index, arr) {
          if (child.id === nod.id) {
            arr.splice(index, 1);
            if (window.tree.has(nod.id)) {
              window.tree.delete(nod.id);
            }
            console.info("remove node", nod.id)
          }
        });
      });
    });
    */
  }

  link = (parentid, childid) => {

    let onods = this.state.nods
    let children

    console.info("link", parentid, "<=", childid)

    for (var i = 0; i < onods.length; i++) {

      if (onods[i].id == childid) {
        children = onods[i]
        onods.splice(i, 1)
        break
      }

      this.findNode(childid, onods[i], (parent, innerChildren, idx) => {

        parent.children.splice(idx, 1)
        children = innerChildren

      })
    }

    console.info("link children", children)
    if (children) {
      for (var i = 0; i < onods.length; i++) {

        if (onods[i].id === parentid) {
          onods[i].children.push(children)
          break
        }

        this.findNode(parentid, onods[i], (_, parent) => {
          parent.children.push(children)
        })
      }
    }

    console.info("link", onods)
    this.setState({ nods: onods })
  }

  unLink = (childid) => {
    let onods = this.state.nods
    let children

    for (var i = 0; i < onods.length; i++) {

      if (onods[i].id == childid) {
        children = onods[i]
        onods.splice(i, 1)
        break
      }

      this.findNode(childid, onods[i], (innerParent, innerChildren, idx) => {
        innerParent.children.splice(idx, 1)
        children = innerChildren
      })

    }

    if (children) {
      onods.push(children)
    }

    this.setState({ nods: onods })
  }

  findNode = (id, parent, callback) => {

    if (parent.children && parent.children.length) {

      for (var i = 0; i < parent.children.length; i++) {

        if (parent.children[i].id == id) {
          callback(parent, parent.children[i], i)
          break
        }

        this.findNode(id, parent.children[i], callback)

      }

    }
  };

  fillData(org, info, graph, edit) {

    console.info("fill", info)

    if (graph) {
      org.pos = info.pos
    }

    if (edit) {
      if (info.ty === NodeTy.Action) {
        org.code = info.code
        org.alias = info.alias
      } else if (info.ty === NodeTy.Assert || info.ty === NodeTy.Condition) {
        org.code = info.code
      } else if (info.ty === NodeTy.Loop) {
        org.loop = info.loop
      } else if (info.ty === NodeTy.Wait) {
        org.wait = info.wait
      }
    }

    org.ty = info.ty

  }

  updateGraphInfo(graphinfo) {

    let tnode = window.tree.get(graphinfo.id)
    this.fillData(tnode, graphinfo, true, false)

    window.tree.set(tnode.id, tnode)

  }

  updateEditInfo(editinfo, notify) {

    let tnode = window.tree.get(editinfo.id)

    this.fillData(tnode, editinfo, false, true)

    if (notify) {
      message.success("apply info succ");
    }

    window.tree.set(editinfo.id, tnode);
  }

  foreachRelation(parent) {

    for (var i = 0; i < parent.children.length; i++) {

      if (window.tree.has(parent.children[i].id)) {
        this.fillData(parent.children[i], window.tree.get(parent.children[i].id), true, true)
      }

      if (parent.children[i].children && parent.children[i].children.length) {
        this.foreachRelation(parent.children[i])
      }

    }

  }

  getTree() {

    let root
    for (var i = 0; i < this.state.nods.length; i++) {
      if (this.state.nods[i].id === this.state.rootid) {
        root = this.state.nods[i]
        break
      }
    }

    this.fillData(root, window.tree.get(root.id), true, false)
    if (root && root.children.length) {
      this.foreachRelation(root)
    }

    return root
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

    PubSub.subscribe(Topic.NodeAdd, (topic, nodeinfo) => {
      this.addNode(nodeinfo);
    });

    PubSub.subscribe(Topic.NodeRmv, (topic, nodeid) => {
      this.rmvNode(nodeid);
    });

    PubSub.subscribe(Topic.LinkConnect, (topic, linkinfo) => {
      this.link(linkinfo.parent, linkinfo.child);
    });

    PubSub.subscribe(Topic.LinkDisconnect, (topic, nodeid) => {
      this.unLink(nodeid);
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

      console.info("debug", tree)

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
            PubSub.publish(Topic.UpdateBlackboard, json.Body.Blackboard);
          } else {

            let changestr, metastr
            let meta = JSON.parse(json.Body.Blackboard)
            let change = JSON.parse(json.Body.Change)

            metastr = JSON.stringify(meta)
            changestr = JSON.stringify(change)

            PubSub.publish(Topic.UpdateBlackboard, metastr);
            PubSub.publish(Topic.UpdateChange, changestr)
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

  componentDidMount() { }

  render() {
    return <div></div>;
  }
}
