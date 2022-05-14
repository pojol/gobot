import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message } from "antd";

import Topic from "../../../constant/topic";

export default class Blackboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      metadata: {},
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
    });
  }

  render() {
    return (
      <div>
        <ReactJson
          name="Meta"
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
