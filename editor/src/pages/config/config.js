import {
  Input,
  Button,
  message,
  Select,
  Collapse,
  InputNumber,
} from "antd";
import * as React from "react";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/theme/abcdef.css";
import "codemirror/theme/ayu-dark.css";
import "codemirror/theme/yonce.css";
import "codemirror/theme/neo.css";
import "codemirror/theme/zenburn.css";
import "codemirror/mode/lua/lua";

import moment from "moment";
import lanMap from "../../locales/lan";
import Api from "../../constant/api";
import { PostBlob, PostGetBlob, CheckHealth, Post } from "../../utils/request";

const { Search } = Input;
const { Option } = Select;
const { Panel } = Collapse;

export default class BotConfig extends React.Component {
  newTabIndex = 0;

  constructor(props) {
    super(props);
    this.state = {
      driveAddr: "",
      globalPrefab: "",
      theme: "ayu-dark",
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
    console.info("syncConfig ======>")
    Post(localStorage.remoteAddr, Api.ConfigSystemInfo, {}).then((json) => {
      if (json.Code !== 200) {
        message.error(
          "get system config fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        console.info(json.Body)
        this.setState({ "reportsize": json.Body.ReportSize, "channelsize": json.Body.ChannelSize })
      }
    })

    PostGetBlob(localStorage.remoteAddr, Api.ConfigGlobalInfo, {}).then((file) => {

      let callback = (content) => {
        this.setState({ globalPrefab: content })
      }

      let reader = new FileReader();
      reader.onload = function (ev) {
        callback(reader.result)
      }
      reader.readAsText(file.blob);
    }

    )};

  onApplyDriveAddr = () => {
    if (this.isUrl(this.state.driveAddr)) {
      let driveAddr = this.state.driveAddr;

      CheckHealth(driveAddr).then((res) => {
        console.info("check health", driveAddr, res);
        if (res.code !== 200) {
          message.error("server connection error " + res.code.toString());
        } else {
          // reset
          window.prefab = new Map();
          localStorage.remoteAddr = driveAddr;
          this.syncConfig();
        }
      });
    } else {
      message.warning("Please enter a valid address");
    }
  };

  onApplyCode = () => {
    
    var blob = new Blob([this.state.globalPrefab], {
      type: "application/json",
    });

    PostBlob(
      localStorage.remoteAddr,
      Api.ConfigGlobalSet,
      "global",
      blob
    ).then((json) => {
      if (json.Code !== 200) {
        message.error(
          "upload fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        message.success("upload succ ");
      }
    });

  };

  onBeforeChange = (editor, data, value) => {
    this.setState({ globalPrefab: value });
  };

  onTableChange = (activeKey) => {
    this.setState({ activeKey });
  };

  onTableEdit = (targetKey, action) => {
    console.info(action, targetKey);
    this[action](targetKey);
  };

  clickTheme = (e) => {
    this.setState({ theme: e }, () => {
      localStorage.theme = e;
    });
  };

  changeChannelSize = (val) => {
    this.setState({ channelsize: val })
  }

  onClickSubmit = () => {
    
    Post(localStorage.remoteAddr, Api.ConfigSystemSet, {
      "ChannelSize": this.state.channelsize,
      "ReportSize": this.state.reportsize,
    }).then((json)=>{
      if (json.Code !== 200) {
        message.error(
          "set config fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        console.info(json.Body)
        message.success("upload succ ");
      }
    })

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
            <CodeMirror
              value={this.state.globalPrefab}
              options={options}
              onBeforeChange={this.onBeforeChange}
            />
            <Button type="primary" onClick={this.onApplyCode}>
              {lanMap["app.config.code.apply"][moment.locale()]}
            </Button>
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
      </div>
    );
  }
}
