import { Input, Tag, Divider, Button } from "antd";
import * as React from "react";
import PubSub from "pubsub-js";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";
import Topic from "../model/topic";
import Config from "../model/config";

const { Search } = Input;

export default class BotConfig extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      driveAddr: "",
      httpActionCode: Config.httpCode,
      assertCode: Config.assertCode,
      conditionCode: Config.conditionCode,
    };
  }

  componentDidMount() {
    var remote = localStorage.remoteAddr
    if (remote !== undefined && remote !== "") {
      this.setState({driveAddr:remote})
    } else {
      this.setState({driveAddr:Config.driveAddr})
    }
  }

  onChangeDriveAddr = (e) => {
    this.setState({ driveAddr: e.target.value });
  };

  onApplyDriveAddr = () => {
    PubSub.publish(Topic.ConfigUpdate, {
        key: "addr",
        val: this.state.driveAddr,
      });
  };

  onBeforeHttpChange = (editor, data, value) => {
    this.setState({ httpActionCode: value });
  };

  onApplyHttpCode = () => {
    PubSub.publish(Topic.ConfigUpdate, {
        key: "httpCode",
        val: this.state.httpActionCode,
      });
  };

  onBeforeAssertChange = (editor, data, value) => {
    this.setState({ assertCode: value });
  };

  onApplyAssertCode = () => {
    PubSub.publish(Topic.ConfigUpdate, {
        key: "assertCode",
        val: this.state.assertCode,
      });
  };

  onBeforeConditionChange = (editor, data, value) => {
    this.setState({ conditionCode: value });
  };

  onApplyConditionCode = () => {
    PubSub.publish(Topic.ConfigUpdate, {
        key: "conditionCode",
        val: this.state.conditionCode,
      });
  };

  render() {
    const addr = this.state.driveAddr;
    const httpCode = this.state.httpActionCode;
    const assertCode = this.state.assertCode;
    const conditionCode = this.state.conditionCode;
    const options = {
      mode: "text/x-lua",
      theme: "solarized dark",
      lineNumbers: true,
    };

    return (
      <div>
        <Divider>
          {" "}
          <Tag color="#2db7f5">drive service address</Tag>
        </Divider>
        <Search
          placeholder={addr}
          onChange={this.onChangeDriveAddr}
          enterButton="Apply"
          onSearch={this.onApplyDriveAddr}
        />
        <Divider>
          {" "}
          <Tag color="#2db7f5">Script</Tag> action code init template
        </Divider>

        <CodeMirror
          width="600"
          height="400"
          value={httpCode}
          options={options}
          onBeforeChange={this.onBeforeHttpChange}
        />
        <Button onClick={this.onApplyHttpCode}>Apply</Button>

        <Divider>
          {" "}
          <Tag color="#2db7f5">Assert</Tag> action code init template
        </Divider>
        <CodeMirror
          width="600"
          height="400"
          value={assertCode}
          options={options}
          onBeforeChange={this.onBeforeAssertChange}
        />
        <Button onClick={this.onApplyAssertCode}>Apply</Button>

        <Divider>
          {" "}
          <Tag color="#2db7f5">Condition</Tag> action code init template
        </Divider>
        <CodeMirror
          width="600"
          height="400"
          value={conditionCode}
          options={options}
          onBeforeChange={this.onBeforeConditionChange}
        />
        <Button onClick={this.onApplyConditionCode}>Apply</Button>
      </div>
    );
  }
}
