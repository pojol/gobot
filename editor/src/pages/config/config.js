import { Input, Divider, Button, Tabs, message, Select, Space } from "antd";
import * as React from "react";
import PubSub from "pubsub-js";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/theme/abcdef.css";
import "codemirror/theme/ayu-dark.css";
import "codemirror/theme/yonce.css";
import "codemirror/theme/neo.css";
import "codemirror/theme/zenburn.css";
import "codemirror/mode/lua/lua";

import Topic from "../../constant/topic";
import moment from 'moment';
import lanMap from "../../locales/lan";
import Api from "../../constant/api";
import { PostBlob, PostGetBlob } from "../../utils/request";


const { Search } = Input;
const { TabPane } = Tabs;
const { Option } = Select;

export default class BotConfig extends React.Component {

  newTabIndex = 0;

  constructor(props) {
    super(props);
    this.state = {
      driveAddr: "",
      activeKey: 'global',
      panes: [],
      theme: "ayu-dark",
    };
  }

  componentDidMount() {
    var remote = localStorage.remoteAddr
    if (remote !== undefined && remote !== "") {
      this.setState({ driveAddr: remote })
    }

    var temp = localStorage.CodeTemplate
    if (temp !== undefined && temp !== "") {
      this.setState({ panes: JSON.parse(temp) })
    }

    PubSub.subscribe(Topic.ConfigUpdate, (topic, info) => {
      if (info.key === "code" && info.val !== "") {
        console.info(info.val)
        this.setState({ panes: JSON.parse(info.val) })
      }
    })

  }

  isUrl(url) {
    var strRegex = '^((https|http|ftp|rtsp|mms)?://)'
      + '?(([0-9a-z_!~*\'().&=+$%-]+: )?[0-9a-z_!~*\'().&=+$%-]+@)?' //ftp的user@ 
      + '(([0-9]{1,3}.){3}[0-9]{1,3}' // IP形式的URL- 199.194.52.184 
      + '|' // 允许IP和DOMAIN（域名） 
      + '([0-9a-z_!~*\'()-]+.)*' // 域名- www. 
      + '([0-9a-z][0-9a-z-]{0,61})?[0-9a-z].' // 二级域名 
      + '[a-z]{2,6})' // first level domain- .com or .museum 
      + '(:[0-9]{1,4})?' // 端口- :80 
      + '((/?)|' // a slash isn't required if there is no file name 
      + '(/[0-9a-z_!~*\'().;?:@&=+$,%#-]+)+/?)$';
    var re = new RegExp(strRegex);
    //re.test() 
    if (re.test(url)) {
      return (true);
    } else {
      return (false);
    }

  }

  onChangeDriveAddr = (e) => {
    this.setState({ driveAddr: e.target.value });
  };

  onApplyDriveAddr = () => {

    if (this.isUrl(this.state.driveAddr)) {

      let driveAddr = this.state.driveAddr
      PostGetBlob(driveAddr, Api.ConfigGet, {}).then((file) => {
        let reader = new FileReader();
        reader.onload = function (ev) {
          localStorage.CodeTemplate = reader.result;
          localStorage.remoteAddr = driveAddr;

          message.success("apply addr succ")
          PubSub.publish(Topic.ConfigUpdate, { key : "code", val: reader.result})
        };
        reader.readAsText(file.blob);
      });

    } else {
      message.warning("Please enter a valid address")
    }

  };

  onApplyCode = () => {

    let panes = this.state.panes
    var templatecode = JSON.stringify(panes)
    var blob = new Blob([templatecode], {
      type: "application/json",
    });

    PostBlob(localStorage.remoteAddr, Api.ConfigUpload, "config", blob).then(
      (json) => {
        if (json.Code !== 200) {
          message.error(
            "upload fail:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          PubSub.publish(Topic.ConfigUpdate, {
            key: "code",
            val: templatecode,
          })
          message.success("upload succ ");
        }
      }
    )
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

  clickTheme = (e) => {
    this.setState({ theme: e }, () => {
      localStorage.theme = e
    })
  }

  render() {
    const addr = this.state.driveAddr;
    const options = {
      mode: "text/x-lua",
      theme: this.state.theme,
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
          <Space>
            {lanMap["app.config.template"][moment.locale()]}
            <Select placeholder={lanMap["app.config.theme"][moment.locale()]}
              style={{ width: 180 }} onChange={this.clickTheme}>
              <Option value="default">default</Option>
              <Option value="abcdef">abcdef</Option>
              <Option value="ayu-dark">ayu-dark</Option>
              <Option value="yonce">yonce</Option>
              <Option value="neo">neo</Option>
              <Option value="solarized dark">solarized dark</Option>
              <Option value="solarized light">solarized light</Option>
              <Option value="zenburn">zenburn</Option>
            </Select>
          </Space>
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
