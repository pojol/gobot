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

    PubSub.subscribe(Topic.Create, (topic, info) => {
      this.setState({ metadata: JSON.parse("{}") });
    })
  }


  render() {
    const isModalVisible = this.state.isModalVisible;

    return (
      <div className="offset">
        <ReactJson
          name={lanMap["app.edit.blackboardJson"][moment.locale()]}
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
