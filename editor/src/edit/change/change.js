import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message, Divider, Button, Tooltip, Modal, Input } from "antd";
import {
  CloudUploadOutlined,
} from "@ant-design/icons";

import Topic from "../../model/topic";
import "./change.css";

export default class ChangeView extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      metadata: {
      },
      behaviorName: "",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.UpdateChange, (topic, info) => {
      try {
        var blackboard = JSON.parse(info);
        this.setState({ metadata: blackboard });
      } catch (err) {
        message.warning("blackboard parse info err");
      }
    });

    PubSub.subscribe(Topic.Upload, (topic, info) => {
      this.setState({ metadata: JSON.parse("{}") });
    });
  }

  showModal = () => {
    this.setState({ isModalVisible: true });
  };

  modalHandleOk = () => {
    this.setState({ isModalVisible: false });
    if (this.state.behaviorName !== "") {
      PubSub.publish(Topic.Upload, this.state.behaviorName);
    } else {
      message.warning("please enter the file name of the behavior tree");
    }
  };

  behaviorNameChange = (e) => {
    this.setState({ behaviorName: e.target.value });
  };

  modalHandleCancel = () => {
    this.setState({ isModalVisible: false });
  };

  debugClick = () => {
    PubSub.publish(Topic.Run, "");
  };

  createClick = () => {
    PubSub.publish(Topic.Create, "");
    this.setState({ metadata: JSON.parse("{}") });
  };

  stepClick = () => {
    PubSub.publish(Topic.Step, "");
  };

  uploadClick = () => {
    this.setState({ isModalVisible: true });
  };


  render() {
    const isModalVisible = this.state.isModalVisible;

    return (
      <div>
        <ReactJson
          name="Step change"
          src={this.state.metadata}
          theme={"rjv-default"}
          enableClipboard={false}
          displayDataTypes={false}
          edit={false}
          add={false}
        ></ReactJson>
      </div>
    );
  }
}
