import React from "react";
import { Tabs } from "antd";
import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";
import ActionTab from "./edit_action";
import LoopTab from "./edit_loop";
import WaitTab from "./edit_wait";
import { NodeTy } from "../../../constant/node_type";

import moment from "moment";
import lanMap from "../../../locales/lan";
import SequenceTab from "./edit_sequence";

const { TabPane } = Tabs;

function GetPane(props) {
  const nodety = props.nodety;
  const dimensions = props.dimensions;

  if (
    nodety === NodeTy.Action ||
    nodety === NodeTy.Condition ||
    nodety === NodeTy.Assert
  ) {
    return <ActionTab dimensions = {dimensions}/>;
  } else if (nodety === NodeTy.Sequence || nodety === NodeTy.Selector) {
    return <SequenceTab />;
  } else if (nodety === NodeTy.Wait) {
    return <WaitTab />;
  } else if (nodety === NodeTy.Loop) {
    return <LoopTab />;
  }

  return <SequenceTab/>
}

export default class Edit extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      nodety: NodeTy.Action,
      nodeid: "",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeClick, (topic, dat) => {
      if (this.state.nodeid !== dat.id) {
        this.setState({ nodety: dat.type, nodeid: dat.id }, ()=>{
          PubSub.publish(Topic.NodeEditorClick, dat);
        });
      }
    });
  }

  render() {
    const { width, height } = this.props.dimensions;

    return (
      <div>
        <GetPane nodety={this.state.nodety} dimensions={this.props.dimensions}/>
      </div>
    );
  }
}
