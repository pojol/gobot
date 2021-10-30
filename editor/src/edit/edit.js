import React, { StrictMode } from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message, Divider, Button, Tooltip, Modal, Input } from "antd";
import {
  CloudUploadOutlined,
  VerticalAlignBottomOutlined,
} from "@ant-design/icons";
import SplitPane, { Pane } from "react-split-pane";

import GraphView from "./graph/graph";
import Edit from "./sider/edit_tab";
import Blackboard from "./sider/blackboard";

import Topic from "../model/topic";

import "./edit.css";

export default class EditPlane extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  componentDidMount() {}

  onEditDragFinished = (e) => {
    console.info("edit", e);
  };

  onCodeDragFinished = (e) => {
    PubSub.publish(Topic.EditPlaneResize, e);
  };

  render() {
    return (
      <div>
        <SplitPane
          split="vertical"
          minSize={500}
          defaultSize="60%"
          minSize={400}
          onDragFinished={this.onEditDragFinished}
        >
          <GraphView />

          <SplitPane
            split="horizontal"
            defaultSize={500}
            minSize={100}
            onDragFinished={this.onCodeDragFinished}
          >
            <Pane minSize={200} maxSize={1000} defaultSize="50%">
              <Edit />
            </Pane>
            <Pane minSize={200} maxSize={1000} defaultSize="50%">
              <Blackboard />
            </Pane>
          </SplitPane>
        </SplitPane>
      </div>
    );
  }
}
