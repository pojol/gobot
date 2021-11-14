import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message, Divider, Button, Tooltip, Modal, Input } from "antd";
import {
  CloudUploadOutlined,
} from "@ant-design/icons";

import Topic from "../../model/topic";
import "./blackboard.css";

export default class Blackboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      metadata: {
      },
      behaviorName: "",

      isModalVisible: false,
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.Blackboard, (topic, info) => {
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
      <div className="offset">
        <Tooltip
          placement="topLeft"
          title="Create a bot based on the current behavior tree"
        >
          <Button onClick={this.createClick}>Create</Button>
        </Tooltip>

        <Divider type="vertical"></Divider>
        <Tooltip
          placement="topLeft"
          title="View the runtime of the behavior tree"
        >
          <Button onClick={this.stepClick}>Step</Button>
        </Tooltip>
        <Divider type="vertical"></Divider>
        <Divider type="vertical"></Divider>
        <Divider type="vertical"></Divider>
        <Tooltip
          placement="topLeft"
          title="Upload the behavior tree file to the server"
        >
          <Button icon={<CloudUploadOutlined />} onClick={this.uploadClick}>
            Upload
          </Button>
        </Tooltip>

        <Modal
          visible={isModalVisible}
          onOk={this.modalHandleOk}
          onCancel={this.modalHandleCancel}
        >
          <Input
            placeholder="input behavior file name"
            onChange={this.behaviorNameChange}
          />
        </Modal>
        <div>
          <ReactJson
            src={this.state.metadata}
            theme={"rjv-default"}
            enableClipboard={false}
            displayDataTypes={false}
            edit={false}
            add={false}
          ></ReactJson>
        </div>
      </div>
    );
  }
}
