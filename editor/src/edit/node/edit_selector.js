import React from "react";
import { Input } from "antd";


export default class SelectorTab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      node_id: "",
      node_ty: "SelectorNode",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeEditorClick, (topic, dat) => {
      this.setState({ node_id: dat.id });
    });
  }

  render() {
    return (
      <div>
        <Input placeholder="id" value={this.state.node_id} disabled={true} />
      </div>
    );
  }
}