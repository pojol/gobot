import React from "react";
import { Input } from "antd";

import PubSub from "pubsub-js";
import Topic from "../../model/topic";

import { Button } from "antd";

export default class SequenceTab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      node_id: "",
      node_ty: "SequenceNode",
    };
  }


  componentDidMount() {
    PubSub.subscribe(Topic.NodeEditorClick, (topic, dat) => {
      this.setState({ node_id: dat.id })
    })
  }

  render() {
    return (
      <div>
        <Button type="dashed">{this.state.node_id}</Button>
      </div>
    )
  }
}