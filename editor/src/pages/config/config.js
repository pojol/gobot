import {
  Input,
  Divider,
  Button,
  Tabs,
  message,
  Select,
  Space,
  Modal,
  Switch,
} from "antd";
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
import { SwatchesPicker } from 'react-color';


import Topic from "../../constant/topic";
import moment from "moment";
import lanMap from "../../locales/lan";
import Api from "../../constant/api";
import { PostBlob, PostGetBlob, CheckHealth, Post } from "../../utils/request";

const { Search } = Input;
const { TabPane } = Tabs;
const { Option } = Select;

export default class BotConfig extends React.Component {
  newTabIndex = 0;

  constructor(props) {
    super(props);
    this.state = {
      driveAddr: "",
      activeKey: "Global",
      panes: [],
      theme: "ayu-dark",
      isModalVisible: false,
      modalConfig: "",
      switchChecked: true,
    };
  }

  componentDidMount() {
    var remote = localStorage.remoteAddr;
    if (remote !== undefined && remote !== "") {
      this.setState({ driveAddr: remote });
    }

    let configmap = window.config;
    var oldpanes = this.state.panes;

    configmap.forEach(function (value, key) {
      var jobj = JSON.parse(value);
      oldpanes.push({
        title: key,
        content: jobj["content"],
        key: key,
        closable: jobj["closable"],
        prefab: jobj["prefab"],
        color: "#fff",
      });
    });
    this.setState({ panes: oldpanes });

    PubSub.subscribe(Topic.ConfigUpdate, (topic, info) => {

      var oldpanes = this.state.panes;
      console.info("info", info)
      var jobj = JSON.parse(info)

      oldpanes.push({
        title: jobj["title"],
        content: jobj["content"],
        key: jobj["title"],
        closable: jobj["closable"],
        prefab: jobj["prefab"],
      });

      this.setState({ panes: oldpanes });

    });
  }

  appendPane(val) {

  }

  isUrl(url) {
    var strRegex =
      "^((https|http|ftp|rtsp|mms)?://)" +
      "?(([0-9a-z_!~*'().&=+$%-]+: )?[0-9a-z_!~*'().&=+$%-]+@)?" + //ftp的user@
      "(([0-9]{1,3}.){3}[0-9]{1,3}" + // IP形式的URL- 199.194.52.184
      "|" + // 允许IP和DOMAIN（域名）
      "([0-9a-z_!~*'()-]+.)*" + // 域名- www.
      "([0-9a-z][0-9a-z-]{0,61})?[0-9a-z]." + // 二级域名
      "[a-z]{2,6})" + // first level domain- .com or .museum
      "(:[0-9]{1,4})?" + // 端口- :80
      "((/?)|" + // a slash isn't required if there is no file name
      "(/[0-9a-z_!~*'().;?:@&=+$,%#-]+)+/?)$";
    var re = new RegExp(strRegex);
    //re.test()
    if (re.test(url)) {
      return true;
    } else {
      return false;
    }
  }

  onChangeDriveAddr = (e) => {
    this.setState({ driveAddr: e.target.value });
  };

  syncConfig = () => {
    console.info("sync config", localStorage.remoteAddr + "/" + Api.ConfigList)
    Post(localStorage.remoteAddr, Api.ConfigList, {}).then((json) => {
      if (json.Code !== 200) {
        message.error(
          "get config list fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        let lst = json.Body.Lst;
        console.info("config lst", lst)
        var counter = 0;

        lst.forEach(function (element) {
          PostGetBlob(localStorage.remoteAddr, Api.ConfigGet, element).then(
            (file) => {
              let reader = new FileReader();
              reader.onload = function (ev) {
                window.config.set(element, reader.result)

                PubSub.publish(Topic.ConfigUpdate, reader.result);

                counter++
                if (counter === lst.length) {
                  PubSub.publish(Topic.ConfigUpdateAll, {})
                }
              };

              reader.readAsText(file.blob);
            }
          );
        });

      }
    });
  }

  onApplyDriveAddr = () => {
    if (this.isUrl(this.state.driveAddr)) {
      let driveAddr = this.state.driveAddr;

      CheckHealth(driveAddr).then((res) => {
        console.info("check health", res);
        if (res.code !== 200) {
          message.error("server connection error " + res.code.toString());
        } else {
          // reset
          window.config = new Map();
          this.setState({ panes: [] })
          this.syncConfig()
        }
      });
    } else {
      message.warning("Please enter a valid address");
    }
  };

  onApplyCode = () => {

    const { panes, activeKey } = this.state;
    const selectPanes = panes.filter((pane) => pane.key === activeKey);

    if (selectPanes.length) {
      let selectPane = selectPanes[0]

      var templatecode = JSON.stringify(selectPane);
      var blob = new Blob([templatecode], {
        type: "application/json",
      });

      console.info("apply config code", templatecode)

      PostBlob(
        localStorage.remoteAddr,
        Api.ConfigUpload,
        selectPane["title"],
        blob
      ).then((json) => {
        if (json.Code !== 200) {
          message.error(
            "upload fail:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          message.success("upload succ ");
          window.config.set(activeKey, templatecode)
          PubSub.publish(Topic.ConfigUpdateAll, {}) // reload
        }
      });

    }


  };

  onBeforeChange = (editor, data, value) => {
    console.info(this.state.activeKey, value);
    let activeKey = this.state.activeKey;

    let newPanes = this.state.panes;
    for (var i = 0; i < newPanes.length; i++) {
      if (newPanes[i].key === activeKey) {
        newPanes[i].content = value;
      }
    }

    this.setState({ panes: newPanes });
  };

  onTableChange = (activeKey) => {
    this.setState({ activeKey });
  };

  onTableEdit = (targetKey, action) => {
    console.info(action, targetKey);
    this[action](targetKey);
  };

  add = () => {
    this.setState({ modalConfig: "" })
    this.showModal();

  };

  remove = (targetKey) => {

    const { panes, activeKey } = this.state;
    let newActiveKey = activeKey;
    let lastIndex;
    let find
    panes.forEach((pane, i) => {
      if (pane.key === targetKey) {
        lastIndex = i - 1;
        find = true
      }
    });

    if (!find) {
      message.warning("unknow template : " + targetKey)
      return
    }

    const newPanes = panes.filter((pane) => pane.key !== targetKey);
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
    }, () => {

      Post(localStorage.remoteAddr, Api.ConfigRemove, { Name: targetKey }).then((json) => {
        if (json.Code !== 200) {
          message.error(
            "remove template err:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          window.config.delete(targetKey)
          PubSub.publish(Topic.ConfigUpdateAll, {}) // reload
          message.success("remove template " + targetKey + " succ")
        }
      });

    });
  };

  clickTheme = (e) => {
    this.setState({ theme: e }, () => {
      localStorage.theme = e;
    });
  };

  showModal = () => {
    this.setState({ isModalVisible: true });
  };

  modalConfigChange = (e) => {
    this.setState({ modalConfig: e.target.value });
  };

  modalHandleOk = () => {
    this.setState({ isModalVisible: false });
    if (this.state.modalConfig !== "") {
      const { panes, modalConfig, switchChecked } = this.state;
      var repeat = false;

      panes.forEach(function (value) {
        if (value.title === modalConfig) {
          message.warning("Duplicate titles cannot be used!");
          repeat = true;
          return;
        }
      });

      if (repeat) {
        return;
      }

      const newPanes = [...panes];
      newPanes.push({
        title: modalConfig,
        content: "Content of new Tab",
        key: modalConfig,
        prefab: switchChecked,
      });
      this.setState({ panes: newPanes, activeKey: modalConfig, modalConfig: "" });
    }
  };

  modalHandleCancel = () => {
    this.setState({ isModalVisible: false });
  };

  switchChange = (checked) => {
    this.setState({ switchChecked: checked })
  }

  handleColorChange = (color) => {
    console.info(color.hex)
    this.setState({color: color.hex})
  }

  handleOnRemove = (value) => {
    
    this.remove(value)

  }

  render() {
    const addr = this.state.driveAddr;
    const options = {
      mode: "text/x-lua",
      theme: this.state.theme,
      lineNumbers: true,
    };

    const { panes, activeKey, isModalVisible } = this.state;

    return (
      <div>
        <Divider>{lanMap["app.config.drive.address"][moment.locale()]}</Divider>
        <Search
          placeholder={addr}
          onChange={this.onChangeDriveAddr}
          enterButton={lanMap["app.config.drive.apply"][moment.locale()]}
          onSearch={this.onApplyDriveAddr}
        />
        <Divider>
          <Space>
            {lanMap["app.config.template"][moment.locale()]}
            <Select
              placeholder={lanMap["app.config.theme"][moment.locale()]}
              style={{ width: 180 }}
              onChange={this.clickTheme}
            >
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
          {panes.map((pane) => (
            <TabPane tab={pane.title} key={pane.key} closable={false}>
              <CodeMirror
                value={pane.content}
                options={options}
                onBeforeChange={this.onBeforeChange}
              />
            </TabPane>
          ))}
        </Tabs>
        <Button type="primary" onClick={this.onApplyCode}>
          {lanMap["app.config.code.apply"][moment.locale()]}
        </Button>

        <Modal
          visible={isModalVisible}
          onOk={this.modalHandleOk}
          onCancel={this.modalHandleCancel}
        >
          <Input
            placeholder={lanMap["app.config.modal.input"][moment.locale()]}
            onChange={this.modalConfigChange}
            value={this.state.modalConfig}
          />
          <Switch checkedChildren={lanMap["app.config.modal.checked"][moment.locale()]} unCheckedChildren={lanMap["app.config.modal.uncheked"][moment.locale()]} onChange={this.switchChange} defaultChecked />
          <SwatchesPicker color={ this.state.color } onChange={ this.handleColorChange }/>
        </Modal>

    <Search
      placeholder="input remove template name"
      allowClear
      enterButton="Remove template"
      size="large"
      onSearch={this.handleOnRemove}
    />
      </div>
    );
  }
}
