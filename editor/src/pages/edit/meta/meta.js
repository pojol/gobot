import React from "react";
import ReactJson from "react-json-view";
import PubSub from "pubsub-js";
import { message, Tabs } from "antd";
import { CodeOutlined, FileSearchOutlined } from "@ant-design/icons";

import Topic from "../../../constant/topic";
import Editor from "react-medium-editor";

require("medium-editor/dist/css/medium-editor.css");
require("medium-editor/dist/css/themes/default.css");

const { TabPane } = Tabs;

export default class Blackboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      metadata: {},
      context: "",
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

    PubSub.subscribe(Topic.UpdateChange, (topic, info) => {
      try {
        info.msg += "\n\n"
        this.setState({ context: info.msg });
      } catch (err) {
        message.warning("blackboard parse info err");
      }
    });

    PubSub.subscribe(Topic.Upload, (topic, info) => {
      this.setState({ context: "" });
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
        <Tabs defaultActiveKey="1">
          <TabPane
            tab={
              <span>
                <FileSearchOutlined />
                Meta
              </span>
            }
            key="2"
          >
            <ReactJson
              name=""
              src={this.state.metadata}
              theme={"rjv-default"}
              enableClipboard={false}
              displayDataTypes={false}
              edit={false}
              add={false}
            ></ReactJson>
          </TabPane>
          <TabPane
            tab={
              <span>
                <CodeOutlined />
                Response
              </span>
            }
            key="1"
          >
            <Editor
              tag="pre"
              //https://github.com/yabwe/medium-editor/blob/d113a74437fda6f1cbd5f146b0f2c46288b118ea/OPTIONS.md#disableediting
              options={{
                placeholder: { text: "", hideOnClick: true },
                disableEditing: true,
              }}
              text={this.state.context}
            />
          </TabPane>

        </Tabs>
      </div>
    );
  }
}
