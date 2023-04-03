import { Input, Button, message, Select, Collapse, InputNumber } from "antd";
import React, { useState, useEffect } from "react";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/theme/abcdef.css";
import "codemirror/theme/ayu-dark.css";
import "codemirror/theme/yonce.css";
import "codemirror/theme/neo.css";
import "codemirror/theme/zenburn.css";
import "codemirror/mode/lua/lua";

import Api from "../constant/api";
const {
  PostBlob,
  PostGetBlob,
  CheckHealth,
  Post,
} = require("../utils/request");

const { Search } = Input;
const { Option } = Select;
const { Panel } = Collapse;

export default function Config() {
  const [state, setState] = useState({
    driveAddr: localStorage.remoteAddr || "",
    globalPrefab: "",
    theme: "ayu-dark",
    switchChecked: true,
    reportsize: 0,
    channelsize: 0,
  });

  useEffect(() => {
    // 在组件初始化时调用一次
    console.info("state.driveAddr", state.driveAddr);
    if (state.driveAddr != "") {
      syncConfig();
    }
  }, []);

  const isUrl = (url: string): boolean => {
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
  };

  const onChangeDriveAddr = (e: any) => {
    setState((state) => ({ ...state, driveAddr: e.target.value }));
  };

  const syncConfig = () => {
    console.info("syncConfig ======>");
    Post(localStorage.remoteAddr, Api.ConfigSystemInfo, {}).then(
      (json: any) => {
        if (json.Code !== 200) {
          message.error(
            "get system config fail:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          console.info("syncConfig body", json.Body);
          setState((state) => ({
            ...state,
            reportsize: json.Body.ReportSize,
            channelsize: json.Body.ChannelSize,
          }));
        }
      }
    );

    PostGetBlob(localStorage.remoteAddr, Api.ConfigGlobalInfo, {}).then(
      (file: any) => {
        let callback = (content: any) => {
          setState((state) => ({ ...state, globalPrefab: content }));
        };

        let reader = new FileReader();
        reader.onload = function (ev) {
          callback(reader.result);
        };
        reader.readAsText(file.blob);
      }
    );
  };

  const onApplyDriveAddr = () => {
    if (isUrl(state.driveAddr)) {
      let driveAddr = state.driveAddr;

      CheckHealth(driveAddr).then((res: any) => {
        console.info("check health", driveAddr, res);
        if (res.code !== 200) {
          message.error("server connection error " + res.code.toString());
        } else {
          // reset
          //const dispatch = useDispatch();
          //dispatch(cleanItems());

          localStorage.remoteAddr = driveAddr;
          syncConfig();
        }
      });
    } else {
      message.warning("Please enter a valid address");
    }
  };

  const onApplyCode = () => {
    var blob = new Blob([state.globalPrefab], {
      type: "application/json",
    });

    PostBlob(localStorage.remoteAddr, Api.ConfigGlobalSet, "global", blob).then(
      (json: any) => {
        if (json.Code !== 200) {
          message.error(
            "upload fail:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          message.success("upload succ ");
        }
      }
    );
  };

  const onBeforeChange = (editor: any, data: any, value: any) => {
    setState((state) => ({ ...state, globalPrefab: value }));
  };

  const clickTheme = (e: any) => {
    setState((state) => ({ ...state, theme: e }));
    localStorage.codeboxTheme = e;
  };

  const changeChannelSize = (val: any) => {
    setState((state) => ({ ...state, channelsize: val }));
  };

  const onClickSubmit = () => {
    Post(localStorage.remoteAddr, Api.ConfigSystemSet, {
      ChannelSize: state.channelsize,
      ReportSize: state.reportsize,
    }).then((json: any) => {
      if (json.Code !== 200) {
        message.error(
          "set config fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        console.info(json.Body);
        message.success("upload succ ");
      }
    });
  };

  const changeReportSize = (val: any) => {
    setState({
      ...state,
      reportsize: val,
    });
  };

  return (
    <div>
      <Collapse defaultActiveKey={["1", "2", "3", "4", "5"]}>
        <Panel
          header={"Drive server address e.g. http://127.0.0.1:8888"}
          key="1"
        >
          <Search
            placeholder={state.driveAddr}
            onChange={onChangeDriveAddr}
            enterButton={"Apply"}
            onSearch={onApplyDriveAddr}
          />
        </Panel>
        <Panel header={"Select a codebox theme"} key="2">
          <Select style={{ width: 200 }} onChange={clickTheme}>
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
        <Panel header={"Global script node"} key="3">
          <CodeMirror
            value={state.globalPrefab}
            options={{
              mode: "lua",
              theme: localStorage.codeboxTheme,
              lineNumbers: true,
            }}
            onBeforeChange={onBeforeChange}
          />
          <Button type="primary" onClick={onApplyCode}>
            {"Apply"}
          </Button>
        </Panel>
        <Panel header={"The number of concurrent robots"} key="4">
          <Input.Group compact>
            <InputNumber
              style={{
                width: "calc(100% - 200px)",
              }}
              min={1}
              value={state.channelsize}
              onChange={changeChannelSize}
            />
            <Button type="primary" onClick={onClickSubmit}>
              Submit
            </Button>
          </Input.Group>
        </Panel>
        <Panel header={"Number of reports archived"} key="5">
          <Input.Group compact>
            <InputNumber
              style={{
                width: "calc(100% - 200px)",
              }}
              min={1}
              max={10000}
              value={state.reportsize}
              onChange={changeReportSize}
            />
            <Button type="primary" onClick={onClickSubmit}>
              Submit
            </Button>
          </Input.Group>
        </Panel>
      </Collapse>
    </div>
  );
}
