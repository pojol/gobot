import React from "react";

import { NodeTy } from "../../../constant/node_type";

import ActionTab from "./action";
import LoopTab from "./loop";
import WaitTab from "./wait";
import SequenceTab from "./sequence";

import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";

/// <reference path="node.d.ts" />

function GetPane(state: any) {
  const nodety = state.nodety;

  switch (nodety) {
    case NodeTy.Sequence:
    case NodeTy.Selector:
    case NodeTy.Root:
      return <SequenceTab />;
    case NodeTy.Wait:
      return <WaitTab />;
    case NodeTy.Loop:
      return <LoopTab  />;
    default:
      return <ActionTab  />;
  }
}

export default class Nodes extends React.Component<Props, {}> {
  state = {
    nodety: "",
    nodeid: "",
  };

  componentDidMount() {
    PubSub.subscribe(Topic.NodeGraphClick, (topic: string, dat: any) => {
      if (this.state.nodeid !== dat.id) {
        this.setState({ nodety: dat.type, nodeid: dat.id }, () => {
          console.info("pub editor click", dat.type, dat.id);
          PubSub.publish(Topic.NodeEditorClick, dat);
        });
      }
    });
  }

  render() {
    return (
      <div>
        <GetPane {...this.state} />
      </div>
    );
  }
}
