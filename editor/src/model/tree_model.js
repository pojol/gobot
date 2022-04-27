import React from "react";
import PubSub from "pubsub-js";
import Topic from "./topic";
import { message,Alert } from "antd";
import OBJ2XML from "object-to-xml";
import { Post, PostBlob } from "./request";
import Api from "./api";
import { NodeTy } from "./node_type";


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

const Cmd = {
  ADD: "nod_add",
  RMV: "nod_rmv",
  Update: "nod_update",
  Link: "node_link",
  Unlink: "node_unlink",
}

function ErrMsgParse(msg) {

  var arr = msg.split("\n")
  var newmsg = ""
  var wline,tline = -1

  for (var i = 0; i < arr.length; i++) {
    var lidx = arr[i].indexOf("line:")
    var cidx = arr[i].indexOf("(column:")
    if (lidx !== -1 && cidx !== -1) {
      wline = parseInt(arr[i].substring(lidx+5, cidx), 10) -1
      tline = i
      console.info("err line", wline, "tips", tline)
    }
  }

  for(var i =0; i < arr.length; i++) {
    if (tline == i) {
      newmsg += "<u>" + arr[i] + "</u>\n"
    } else if (wline == i) {
      newmsg += "<font color='red'>" + arr[i] + "</font>\n"
      console.info("warp", newmsg)
    } else {
      newmsg += arr[i]+"\n"
    }
  }

  newmsg += "\n\n"

  return newmsg
}

export default class TreeModel extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      rootid: "",
      nods: [], //  root 记录节点的链路关系， window(map 记录节点的细节
      botid: "",
      behaviorTreeName: "",
      httpCodeTmp: "",
      assertTmp: `
-- Write expression to return true or false
function execute()

end
      `,
      conditionTmp: `
-- Write expression to return true or false
function execute()

end
      `,
      history: [],
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

    if (nod.children && nod.children.length) {
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

  addNode = (nod, silent) => {
    if (nod.ty === NodeTy.Root) {
      this.setState({ rootid: nod.id })
    }

    let rinfo = this.getRelationInfo(nod)
    this.syncMapInfo(nod)

    let olst = this.state.nods
    olst.push(rinfo)

    let ohistory = this.state.history
    if (!silent) {
      let cmd = [{ "cmd": Cmd.RMV, "parm": nod.id }]
      console.info("history push", cmd)
      ohistory.push(cmd)
    }

    this.setState({ nods: olst, history: ohistory })
  }

  rmvNode = (id, silent) => {

    if (id === this.state.rootid) {
      return
    }

    let rmvnod, rmvparent
    let nnods = this.state.nods
    let ohistory = this.state.history

    for (var i = 0; i < nnods.length; i++) {

      if (nnods[i].id === id) {
        rmvnod = nnods[i]
        nnods.splice(i, 1)
        break
      }

      this.findNode(id, nnods[i], (parent, children, idx) => {
        parent.children.splice(idx, 1)
        rmvparent = parent
        rmvnod = children
      })

      if (rmvnod) {
        this.fillData(rmvnod, window.tree.get(rmvnod.id), true, true)
        this.foreachRelation(rmvnod)

        if (!silent) {
          let cmd = [
            { "cmd": Cmd.ADD, "parm": rmvnod },
            { "cmd": Cmd.Link, "parm": [rmvparent.id, rmvnod.id] }
          ]
          console.info("history push", cmd)
          ohistory.push(cmd)
        }

        this.walk(rmvnod, (nod) => {
          if (window.tree.has(nod.id)) {
            window.tree.delete(nod.id);
          }
        })
      }

      this.setState({ nods: nnods, history: ohistory })
    }

  }

  findTree = (nods, id) => {

    for (var i = 0; i < nods.length; i++) {

      if (nods[i].id === id) {
        return nods[i]
      }

      if (nods[i].chlidren) {
        var res = this.findTree(nods[i].children, id)
        if (res) {
          return res
        }
      }

    }


  }

  link = (parentid, childid, silent) => {

    let children
    let onods = this.state.nods
    let ohistory = this.state.history

    for (let i = 0; i < onods.length; i++) {

      if (onods[i].id === childid) {
        children = onods[i]
        onods.splice(i, 1)
        break
      }

      this.findNode(childid, onods[i], (parent, innerChildren, idx) => {

        if (!silent) {
          let cmd = [
            { "cmd": Cmd.unLink, "parm": [innerChildren.id] }
          ]
          console.info("history push", cmd)
          ohistory.push(cmd)
        }

        parent.children.splice(idx, 1)
        children = innerChildren

      })

    }

    if (children) {
      for (let i = 0; i < onods.length; i++) {

        if (onods[i].id === parentid) {
          onods[i].children.push(children)
          break
        }

        this.findNode(parentid, onods[i], (_, parent) => {
          parent.children.push(children)
        })
      }
    }

    this.setState({ nods: onods })

  }

  unLink = (childid, silent) => {
    let onods = this.state.nods
    let children
    let ohistory = this.state.history

    for (var i = 0; i < onods.length; i++) {

      if (onods[i].id === childid) {
        children = onods[i]
        onods.splice(i, 1)
        break
      }

      this.findNode(childid, onods[i], (innerParent, innerChildren, idx) => {

        if (!silent) {
          let cmd = [
            { "cmd": Cmd.Link, "parm": [innerParent.id, innerChildren.id] }
          ]
          console.info("history push", cmd)
          ohistory.push(cmd)
        }

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

        if (parent.children[i].id === id) {
          callback(parent, parent.children[i], i)
          break
        }

        this.findNode(id, parent.children[i], callback)

      }

    }
  };

  walk = (tree, callback) => {

    if (tree.children && tree.children.length) {

      for (var i = 0; i < tree.children.length; i++) {
        callback(tree.children[i])

        this.walk(tree.children[i], callback)
      }

    }

  }

  fillData(org, info, graph, edit) {

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

  getAllTree() {

    let nods = []

    for (var i = 0; i < this.state.nods.length; i++) {
      var nod = this.state.nods[i]
      this.fillData(nod, window.tree.get(nod.id), true, true)

      if (nod.children && nod.children.length) {
        this.foreachRelation(nod)
      }

      nods.push(nod)
    }

    return nods
  }

  undo() {
    let ohistory = this.state.history
    if (ohistory.length) {

      let h = ohistory.pop()
      console.info("history pop", h)

      for (var i = 0; i < h.length; i++) {
        if (h[i].cmd === Cmd.ADD) {
          this.addNode(h[i].parm, true)
        } else if (h[i].cmd === Cmd.RMV) {
          this.rmvNode(h[i].parm, true)
        } else if (h[i].cmd == Cmd.Link) {
          this.link(h[i].parm[0], h[i].parm[1], true)
        } else if (h[i].cmd === Cmd.Unlink) {
          this.Unlink(h[i].parm[0], true)
        }
      }

      let mtree = this.getAllTree()

      //console.info("new tree", JSON.stringify(mtree))

      PubSub.publish(Topic.FileLoadRedraw, mtree)
      this.setState({ history: ohistory })
    }

  }

  componentWillMount() {
    window.tree = new Map(); // 主要维护的是 editor 节点编辑后的数据
    this.setState({ tree: {} }); // 主要维护的是 graph 中节点的数据

    PubSub.subscribe(Topic.ConfigUpdate, (topic, info) => {
      if (info.key === "addr") {
        localStorage.remoteAddr = info.val;
        message.success("addr update succ");
      } else if (info.key === "code" && info.val !== "") {
        var codetmp = JSON.parse(info.val)
        for (var i = 0; i < codetmp.length; i++) {
          if (codetmp[i]["title"] === "HTTP") {
            this.setState({ httpCodeTmp: codetmp[i]["content"] });
          }
        }
      }
    });

    PubSub.subscribe(Topic.NodeAdd, (topic, addinfo) => {
      let info = addinfo[0]
      let build = addinfo[1]
      let silent = addinfo[2]

      if (build) {
        console.info("node model add", info)
        this.addNode(info, silent);
      }
    });

    PubSub.subscribe(Topic.NodeRmv, (topic, nodeid) => {
      this.rmvNode(nodeid);
    });

    PubSub.subscribe(Topic.LinkConnect, (topic, linkinfo) => {

      let info = linkinfo[0]
      let silent = linkinfo[1]

      this.link(info.parent, info.child, silent);
    });

    PubSub.subscribe(Topic.LinkDisconnect, (topic, info) => {
      let nodid = info[0]
      let silent = info[1]

      this.unLink(nodid, silent);
    });

    PubSub.subscribe(Topic.UpdateNodeParm, (topic, info) => {
      this.updateEditInfo(info.parm, info.notify);
    });

    PubSub.subscribe(Topic.UpdateGraphParm, (topic, info) => {
      this.updateGraphInfo(info, false);
    });

    PubSub.subscribe(Topic.Undo, () => {
      this.undo()
    })

    PubSub.subscribe(Topic.HistoryClean, () => {
      console.info("history clean")
      this.setState({ history: [] })
    })

    PubSub.subscribe(Topic.FileLoad, (topic, info) => {
      window.tree = new Map();
      this.setState({ nods: [], rootid: "", behaviorTreeName: info.Name });
    });

    PubSub.subscribe(Topic.Create, (topic, info) => {
      var name = this.state.behaviorTreeName;
      var tree = this.getTree();

      if (name === undefined || name === "") {
        name = tree.id
      }

      var xmltree = {
        behavior: tree,
      };

      var blob = new Blob([OBJ2XML(xmltree)], {
        type: "application/json",
      });

      PostBlob(localStorage.remoteAddr, Api.DebugCreate, name, blob).then(
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

      PostBlob(localStorage.remoteAddr, Api.FileBlobUpload, filename, blob).then(
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

      Post(localStorage.remoteAddr, Api.DebugStep, { BotID: this.state.botid }).then(
        (json) => {
          if (json.Code !== 200) {
            let change
            let changeInfo = {}

            if (json.Code === 1008) {
              change = JSON.parse(json.Body.Change)
              changeInfo = {
                status : "",
                msg : JSON.stringify(change,null,"\t")
              } 

              message.success("the end")
            } else {
              message.error(json.Msg)
              changeInfo = {
                status : "error",
                msg : ErrMsgParse(json.Body.RuntimeErr)
              }
            }
            
            if (json.Code !== 1010) {
              PubSub.publish(Topic.UpdateChange, changeInfo)
            }
            PubSub.publish(Topic.UpdateBlackboard, json.Body.Blackboard);
          } else {

            let metastr
            let meta = JSON.parse(json.Body.Blackboard)
            let change = JSON.parse(json.Body.Change)

            metastr = JSON.stringify(meta)

            PubSub.publish(Topic.UpdateBlackboard, metastr);
            PubSub.publish(Topic.UpdateChange, {status:"", msg:JSON.stringify(change, null, "\t")})
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
