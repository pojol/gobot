import * as React from "react";

import Topic from "../../../constant/topic";
import { NodeTy } from "@/constant/node_type";
import Cmd from "@/constant/cmd";
import { message } from "antd";

import { connect, ConnectedProps } from 'react-redux';
import Api from "@/constant/api";
import store from "@/models/store";



export default class Tree extends React.Component {
  state = {
    nods: new Array<any>(), //  root 记录节点的链路关系， window(map 记录节点的细节
    history: new Array<any>(),
    rootid: "",
    behaviorTreeName: "",
  };


  updateEditInfo(editinfo: NodeNotifyInfo) {
    let tnode = window.tree.get(editinfo.id);
    if (tnode === undefined) {
      tnode = getDefaultNodeNotifyInfo()
    }
    this.fillData(tnode, editinfo, false, true);

    if (editinfo.notify) {
      message.success("apply info succ");
    }

    window.tree.set(editinfo.id, tnode);
  }

  walk = (tree: NodeNotifyInfo, callback: any) => {
    if (tree.children && tree.children.length) {
      for (var i = 0; i < tree.children.length; i++) {
        callback(tree.children[i]);

        this.walk(tree.children[i], callback);
      }
    }
  };

  componentWillMount() {
    console.info("tree model init")
    


    PubSub.subscribe(
      Topic.NodeUpdateParm,
      (topic: string, info: NodeNotifyInfo) => {
        this.updateEditInfo(info);
      }
    );


  }

  render() {
    return <div></div>;
  }
}