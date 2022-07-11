import React from "react";
import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";
import ActionTab from "./edit_action";
import LoopTab from "./edit_loop";
import WaitTab from "./edit_wait";
import { NodeTy } from "../../../constant/node_type";

import SequenceTab from "./edit_sequence";


function GetPane(props) {
  const nodety = props.nodety;
  const dimensions = props.dimensions;

  switch (nodety) {
    case NodeTy.Sequence:
    case NodeTy.Selector:
    case NodeTy.Root:
      return <SequenceTab />;
    case NodeTy.Wait:
      return <WaitTab />;
    case NodeTy.Loop:
      return <LoopTab />;
    default:
      return <ActionTab dimensions={dimensions} />;
  }

}

export default class Edit extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      nodety: "",
      nodeid: "",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeClick, (topic, dat) => {
      if (this.state.nodeid !== dat.id) {
        this.setState({ nodety: dat.type, nodeid: dat.id }, () => {
          PubSub.publish(Topic.NodeEditorClick, dat);
        });
      }
    });
  }

  render() {

    return (
      <div>
        <GetPane nodety={this.state.nodety} dimensions={this.props.dimensions} />
      </div>
    );
  }
}
