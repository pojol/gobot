import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message } from "antd";

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

  stepClick = () => {
    PubSub.publish(Topic.Step, "");
  };


  render() {

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
