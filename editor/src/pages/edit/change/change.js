import React from "react";
import PubSub from "pubsub-js";
import { message } from "antd";

import Topic from "../../../model/topic";
import Editor from 'react-medium-editor';

require('medium-editor/dist/css/medium-editor.css');
require('medium-editor/dist/css/themes/default.css');


export default class ChangeView extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      metadata: "",
      status : "",
      behaviorName: "",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.UpdateChange, (topic, info) => {
      try {
        info.msg += "\n\n"
        console.info(info.msg)
        this.setState({ metadata:info.msg, status: info.status });
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
        <Editor
          tag="pre"
          //https://github.com/yabwe/medium-editor/blob/d113a74437fda6f1cbd5f146b0f2c46288b118ea/OPTIONS.md#disableediting
          options={{ placeholder: { text : "",hideOnClick: true }, disableEditing : true }}
          text={this.state.metadata}
        />
      </div>
    );
  }
}
