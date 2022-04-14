import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { Input, message } from "antd";

import Topic from "../../model/topic";
import "./change.css";

import moment from 'moment';
import lanMap from "../../config/lan";
import { ClockCircleOutlined } from '@ant-design/icons';

const { TextArea } = Input;


export default class ChangeView extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      metadata: "",
      behaviorName: "",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.UpdateChange, (topic, info) => {
      try {
        this.setState({ metadata: info });
      } catch (err) {
        message.warning("blackboard parse info err");
      }
    });

    PubSub.subscribe(Topic.Upload, (topic, info) => {
      this.setState({ metadata: JSON.parse("{}") });
    });
  }


  modalHandleOk = () => {
    if (this.state.behaviorName !== "") {
      PubSub.publish(Topic.Upload, this.state.behaviorName);
    } else {
      message.warning("please enter the file name of the behavior tree");
    }
  };

  behaviorNameChange = (e) => {
    this.setState({ behaviorName: e.target.value });
  };



  debugClick = () => {
    PubSub.publish(Topic.Run, "");
  };

  createClick = () => {
    PubSub.publish(Topic.Create, "");
    this.setState({ metadata: JSON.parse("{}") });
  };

  render() {

    return (
      <div>
        <TextArea
          value={this.state.metadata}
          bordered={false}
          placeholder=""
          disabled={true}
          autoSize={{ minRows: 10 }}
        />
      </div>
    );
  }
}
