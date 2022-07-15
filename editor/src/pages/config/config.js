import {
  Input,
  Button,
  Tabs,
  message,
  Select,
  Modal,
  Switch,
  Collapse,
  InputNumber,
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
import { SwatchesPicker } from "react-color";

import Topic from "../../constant/topic";
import moment from "moment";
import lanMap from "../../locales/lan";
import Api from "../../constant/api";
import { PostBlob, PostGetBlob, CheckHealth, Post } from "../../utils/request";

const { Search } = Input;
const { TabPane } = Tabs;
const { Option } = Select;
const { Panel } = Collapse;

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
      reportsize: 0,
      channelsize: 0,
    };
  }

  componentDidMount() {
    var remote = localStorage.remoteAddr;
    if (remote !== undefined && remote !== "") {
      this.setState({ driveAddr: remote });

      this.syncConfig()
    }

    PubSub.subscribe(Topic.ConfigUpdate, (topic, info) => {
      var oldpanes = this.state.panes;
      var jobj = JSON.parse(info);

      oldpanes.push({
        title: jobj["title"],
        content: jobj["content"],
        key: jobj["title"],
        closable: jobj["closable"],
        prefab: jobj["prefab"],
      });

      this.setState({ panes: oldpanes });
    });

    PubSub.subscribe(Topic.SystemConfigUpdate, (topic, info) => {
      let jobj = JSON.parse(info)

      let reportsize = jobj["reportsize"]
      let channelsize = jobj["channelsize"]

      this.setState({ reportsize: reportsize, channelsize: channelsize })
      console.info("system config", reportsize, channelsize)
    })
  }

  appendPane(val) { }

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
    Post(localStorage.remoteAddr, Api.ConfigList, {}).then((json) => {
      if (json.Code !== 200) {
        message.error(
          "get config list fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        let lst = json.Body.Lst;
        var counter = 0;

        lst.forEach(function (element) {
          PostGetBlob(localStorage.remoteAddr, Api.ConfigGet, element).then(
            (file) => {
              let reader = new FileReader();
              reader.onload = function (ev) {

                console.info("load config", element)
                if (element !== "system") {
                  window.config.set(element, reader.result);
                  PubSub.publish(Topic.ConfigUpdate, reader.result);
                } else {
                  PubSub.publish(Topic.SystemConfigUpdate, reader.result)
                }

                counter++;
                if (counter === lst.length) {
                  PubSub.publish(Topic.ConfigUpdateAll, {});
                }
              };

              reader.readAsText(file.blob);
            }
          );
        });
      }
    });
  };

  onApplyDriveAddr = () => {
    if (this.isUrl(this.state.driveAddr)) {
      let driveAddr = this.state.driveAddr;

      CheckHealth(driveAddr).then((res) => {
        console.info("check health", driveAddr, res);
        if (res.code !== 200) {
          message.error("server connection error " + res.code.toString());
        } else {
          // reset
          window.config = new Map();
          localStorage.remoteAddr = driveAddr;
          this.setState({ panes: [] });
          this.syncConfig();
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
      let selectPane = selectPanes[0];

      var templatecode = JSON.stringify(selectPane);
      var blob = new Blob([templatecode], {
        type: "application/json",
      });

      console.info("apply config code", templatecode);

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
          window.config.set(activeKey, templatecode);
          PubSub.publish(Topic.ConfigUpdateAll, {}); // reload
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
    this.setState({ modalConfig: "" });
    this.showModal();
  };

  remove = (targetKey) => {
    const { panes, activeKey } = this.state;
    let newActiveKey = activeKey;
    let lastIndex;
    let find;
    panes.forEach((pane, i) => {
      if (pane.key === targetKey) {
        lastIndex = i - 1;
        find = true;
      }
    });

    if (!find) {
      message.warning("unknow template : " + targetKey);
      return;
    }

    const newPanes = panes.filter((pane) => pane.key !== targetKey);
    if (newPanes.length && newActiveKey === targetKey) {
      if (lastIndex >= 0) {
        newActiveKey = newPanes[lastIndex].key;
      } else {
        newActiveKey = newPanes[0].key;
      }
    }

    this.setState(
      {
        panes: newPanes,
        activeKey: newActiveKey,
      },
      () => {
        Post(localStorage.remoteAddr, Api.ConfigRemove, {
          Name: targetKey,
        }).then((json) => {
          if (json.Code !== 200) {
            message.error(
              "remove template err:" + String(json.Code) + " msg: " + json.Msg
            );
          } else {
            window.config.delete(targetKey);
            PubSub.publish(Topic.ConfigUpdateAll, {}); // reload
            message.success("remove template " + targetKey + " succ");
          }
        });
      }
    );
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
      var out = false;

      panes.forEach(function (value) {
        if (value.title === modalConfig) {
          message.warning("Duplicate titles cannot be used!");
          out = true;
          throw new Error("break")
        }

        if (modalConfig === "system") {
          message.warning("Can't use keywords `system`");
          out = true;
          throw new Error("break")
        }
      });

      if (out) {
        return;
      }

      const newPanes = [...panes];
      newPanes.push({
        title: modalConfig,
        content: "Content of new Tab",
        key: modalConfig,
        prefab: switchChecked,
      });
      this.setState({
        panes: newPanes,
        activeKey: modalConfig,
        modalConfig: "",
      });
    }
  };

  modalHandleCancel = () => {
    this.setState({ isModalVisible: false });
  };

  switchChange = (checked) => {
    this.setState({ switchChecked: checked });
  };

  handleColorChange = (color) => {
    console.info(color.hex);
    this.setState({ color: color.hex });
  };

  handleOnRemove = (value) => {
    this.remove(value);
  };

  changeChannelSize = (val) => {
    this.setState({ channelsize: val })
  }

  onClickSubmit = () => {

    var templatecode = JSON.stringify({
      "channelsize":this.state.channelsize,
      "reportsize":this.state.reportsize,
    });
    var blob = new Blob([templatecode], {
      type: "application/json",
    });

    PostBlob(
      localStorage.remoteAddr,
      Api.ConfigUpload,
      "system",
      blob
    ).then((json) => {
      if (json.Code !== 200) {
        message.error(
          "upload config fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        message.success("upload succ ");
      }
    });
  }

  changeReportSize = (val) => {
    this.setState({ reportsize: val })
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
        <Collapse defaultActiveKey={["1", "2", "3", "4", "5"]}>
          <Panel
            header={lanMap["app.config.drive.address"][moment.locale()]}
            key="1"
          >
            <Search
              placeholder={addr}
              onChange={this.onChangeDriveAddr}
              enterButton={lanMap["app.config.drive.apply"][moment.locale()]}
              onSearch={this.onApplyDriveAddr}
            />
          </Panel>
          <Panel header={lanMap["app.config.theme"][moment.locale()]} key="2">
            <Select style={{ width: 200 }} onChange={this.clickTheme}>
              <Option value="default">default</Option>
              <Option value="abcdef">abcdef</Option>
              <Option value="ayu-dark">ayu-dark</Option>
              <Option value="yonce">yonce</Option>
              <Option value="neo">neo</Option>
              <Option value="solarized dark">solarized dark</Option>
              <Option value="solarized light">solarized light</Option>
              <Option value="zenburn">zenburn</Option>
            </Select>
          </Panel>
          <Panel
            header={lanMap["app.config.template"][moment.locale()]}
            key="3"
          >
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
            <Search
              placeholder="remove prefab script node by name"
              allowClear
              enterButton="Remove prefab"
              onSearch={this.handleOnRemove}
            />
          </Panel>
          <Panel
            header={lanMap["app.config.channelsize"][moment.locale()]}
            key="4"
          >
            <Input.Group compact>
              <InputNumber
                style={{
                  width: 'calc(100% - 200px)',
                }}
                min={1}
                value={this.state.channelsize}
                onChange={this.changeChannelSize}
              />
              <Button type="primary" onClick={this.onClickSubmit} >Submit</Button>
            </Input.Group>
          </Panel>
          <Panel
            header={lanMap["app.config.reportsize"][moment.locale()]}
            key="5"
          >
            <Input.Group compact>
              <InputNumber
                style={{
                  width: 'calc(100% - 200px)',
                }}
                min={1}
                max={10000}
                value={this.state.reportsize}
                onChange={this.changeReportSize}
              />
              <Button type="primary" onClick={this.onClickSubmit}>Submit</Button>
            </Input.Group>
          </Panel>
        </Collapse>

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
          <Switch
            checkedChildren={
              lanMap["app.config.modal.checked"][moment.locale()]
            }
            unCheckedChildren={
              lanMap["app.config.modal.uncheked"][moment.locale()]
            }
            onChange={this.switchChange}
            defaultChecked
          />
          <SwatchesPicker
            color={this.state.color}
            onChange={this.handleColorChange}
          />
        </Modal>
      </div>
    );
  }
}
