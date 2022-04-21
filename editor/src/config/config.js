import { Input, Tag, Divider, Button, Tabs } from "antd";
import * as React from "react";
import PubSub from "pubsub-js";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";
import Topic from "../model/topic";
import Config from "../model/config";
import moment from 'moment';
import lanMap from "../config/lan";

const { Search } = Input;

const { TabPane } = Tabs;


export default class BotConfig extends React.Component {

  newTabIndex = 0;

  constructor(props) {
    super(props);
    this.state = {
      driveAddr: "",
      activeKey: 'http',
      panes: [
        { title: 'HTTP', content: Config.httpCode, key: 'http', closable: false },
        { title: 'Assert', content: Config.assertCode, key: 'assert', closable: false },
        { title: 'Condition', content: Config.conditionCode, key: 'condition', closable: false },
        { title: "Global", content: Config.globalCode, key: 'global', closable: false},
      ],
    };
  }

  componentDidMount() {
    var remote = localStorage.remoteAddr
    if (remote !== undefined && remote !== "") {
      this.setState({ driveAddr: remote })
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

  onApplyCode = () => {
    PubSub.publish(Topic.ConfigUpdate, {
      key: "assertCode",
      val: this.state.assertCode,
    });
  };

  onBeforeChange = (editor, data, value) => {

    console.info(this.state.activeKey, value)
    let activeKey = this.state.activeKey
    let newPanes = this.state.panes

    for (var i = 0; i < newPanes.length; i++) {
      if (newPanes[i].key === activeKey) {
        newPanes[i].content = value
      }
    }

    this.setState({ panes: newPanes })
  };

  onTableChange = activeKey => {
    this.setState({ activeKey });
  };

  onTableEdit = (targetKey, action) => {
    console.info(action, targetKey)
    this[action](targetKey);
  };

  add = () => {
    const { panes } = this.state;
    const activeKey = `newTab${this.newTabIndex++}`;
    const newPanes = [...panes];
    newPanes.push({ title: 'New Tab', content: 'Content of new Tab', key: activeKey });
    this.setState({
      panes: newPanes,
      activeKey,
    });
  };

  remove = targetKey => {
    const { panes, activeKey } = this.state;
    let newActiveKey = activeKey;
    let lastIndex;
    panes.forEach((pane, i) => {
      if (pane.key === targetKey) {
        lastIndex = i - 1;
      }
    });
    const newPanes = panes.filter(pane => pane.key !== targetKey);
    if (newPanes.length && newActiveKey === targetKey) {
      if (lastIndex >= 0) {
        newActiveKey = newPanes[lastIndex].key;
      } else {
        newActiveKey = newPanes[0].key;
      }
    }
    this.setState({
      panes: newPanes,
      activeKey: newActiveKey,
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

    const { panes, activeKey } = this.state;

    return (
      <div>
        <Divider>
          {lanMap["app.config.drive.address"][moment.locale()]}
        </Divider>
        <Search
          placeholder={addr}
          onChange={this.onChangeDriveAddr}
          enterButton={lanMap["app.config.drive.apply"][moment.locale()]}
          onSearch={this.onApplyDriveAddr}
        />
        <Divider>
          {lanMap["app.config.template"][moment.locale()]}
        </Divider>

        <Tabs
          type="editable-card"
          onChange={this.onTableChange}
          activeKey={activeKey}
          onEdit={this.onTableEdit}
        >
          {panes.map(pane => (
            <TabPane tab={pane.title} key={pane.key} closable={pane.closable}>
              <CodeMirror
                value={pane.content}
                options={options}
                onBeforeChange={this.onBeforeChange}
              />
            </TabPane>
          ))}
        </Tabs>
        <Button type="primary" onClick={this.onApplyCode}>{lanMap["app.config.code.apply"][moment.locale()]}</Button>
      </div>
    );
  }
}
