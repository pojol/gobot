import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message, Button, Tooltip, Modal, Input,Space } from "antd";
import {
  CloudUploadOutlined,
} from "@ant-design/icons";

import moment from 'moment';
import lanMap from "../../config/lan";

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
    PubSub.subscribe(Topic.UpdateBlackboard, (topic, info) => {
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
        <Space direction="vertical">
        <Space>
          <Tooltip
            placement="topLeft"
            title={lanMap["app.edit.blackboard.create.desc"][moment.locale()]}
          >
            <Button onClick={this.createClick}>{lanMap["app.edit.blackboard.create"][moment.locale()]}</Button>
          </Tooltip>
          <Tooltip
            placement="topLeft"
            title={lanMap["app.edit.blackboard.step.desc"][moment.locale()]}
          >
            <Button onClick={this.stepClick}>{lanMap["app.edit.blackboard.step"][moment.locale()]}</Button>
          </Tooltip>
          <Tooltip
            placement="topLeft"
            title={lanMap["app.edit.blackboard.upload.desc"][moment.locale()]}
          >
            <Button icon={<CloudUploadOutlined />} onClick={this.uploadClick}>
            {lanMap["app.edit.blackboard.upload"][moment.locale()]}
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
        </Space>

        <ReactJson
          name={lanMap["app.edit.blackboardJson"][moment.locale()]}
          src={this.state.metadata}
          theme={"rjv-default"}
          enableClipboard={false}
          displayDataTypes={false}
          edit={false}
          add={false}
        ></ReactJson>
        </Space>
        
      </div>
    );
  }
}
