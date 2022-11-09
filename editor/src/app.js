import {
  Layout,
  Tabs,
  Tag,
  Radio,
  Modal,
  Input,
  Image,
  Space,
  message,
  Tooltip,
} from "antd";
import * as React from "react";
import "antd/dist/antd.css";
import "./app.css";
import TreeModel from "./behavior_tree/tree_model";
import PubSub from "pubsub-js";
import Topic from "./constant/topic";

import HomePage from "./pages/home/home";
import ReportPage from "./pages/report/report";
import ConfigPage from "./pages/config/config";
import EditPage from "./pages/edit/edit";
import RunningPage from "./pages/runing/runing";
import BotPrefab from "./pages/prefab/prefab";

import enUS from "antd/lib/locale/en_US";
import zhCN from "antd/lib/locale/zh_CN";
import moment from "moment";
import "moment/locale/zh-cn";
import lanMap from "./locales/lan";
import { Post, PostGetBlob, CheckHealth } from "./utils/request";
import Api from "./constant/api";

import { ReadOutlined, ApiFilled } from "@ant-design/icons";

const { TabPane } = Tabs;
moment.locale("en");

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      tab: "Edit",
      locale: enUS,
      isModalVisible: false,
      modalConfig: "",
      connectColor: "red",
      connectTxt: "Not connected to server, please set in config",
    };
  }

  resizeHandler = () => {
    let resizeTimer;
    if (!resizeTimer) {
      resizeTimer = setTimeout(() => {
        resizeTimer = null;
        PubSub.publish(Topic.WindowResize, {});
      }, 100);
    }
  };

  componentWillMount() {
    if (localStorage.theme === "" || localStorage.theme === undefined) {
      localStorage.theme = "ayu-dark";
    }

    let remote = localStorage.remoteAddr;
    if (remote === "" || remote === undefined) {
      this.setState({ isModalVisible: true });
    } else {
      this.syncTemplateCode();
    }

  }

  componentDidMount() {
    PubSub.subscribe(Topic.FileLoad, (topic, info) => {
      this.setState({ tab: "Edit" });
      PubSub.publish(Topic.FileLoadDraw, [info.Tree]);
    });

    this.checkheath()
    window.addEventListener("resize", this.resizeHandler, false);
  }

  changeTab = (e) => {
    this.setState({ tab: e }, () => {
      if (e === "Home") {
        PubSub.publish(Topic.BotsUpdate, {});
      } else if (e === "Report") {
        PubSub.publish(Topic.ReportUpdate, {});
      } else if (e === "Running") {
        PubSub.publish(Topic.RunningUpdate, {});
      } else if (e === "Edit") {
        this.checkheath()
      }
    });
  };

  checkheath() {
    CheckHealth(localStorage.remoteAddr).then((res => {
      if (res.code === 200) {
        this.setState({ connectColor: "#4caf50", connectTxt: "Connecting" })
      } else {
        this.setState({ connectColor: "red", connectTxt: "Not connected to server, please set in config" })
      }
    }));
  }

  syncTemplateCode() {
    console.info("sync templete config", localStorage.remoteAddr);

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

                if (reader.result.byteLength === 0) {
                  message.warning("get config byte length == 0")
                  return
                }
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

  showModal = () => {
    this.setState({ isModalVisible: true });
  };

  modalConfigChange = (e) => {
    this.setState({ modalConfig: e.target.value });
  };

  modalHandleOk = () => {
    this.setState({ isModalVisible: false });
    if (this.state.modalConfig !== "") {
      localStorage.remoteAddr = this.state.modalConfig;

      this.syncTemplateCode();
      this.checkheath()
    }
  };

  modalHandleCancel = () => {
    this.setState({ isModalVisible: false });
  };

  changeLocale = (e) => {
    const localeValue = e.target.value;
    this.setState({ locale: localeValue });

    moment.locale(localeValue.locale);

    PubSub.publish(Topic.LanuageChange, {});

    console.info("moment=>", moment.locale());
  };

  clickDocTag = () => {
    window.open("https://pojol.gitee.io/gobot/#/");
  };

  clickGithubTag = () => {
    window.open("https://github.com/pojol/gobot");
  };

  render() {
    const { locale, isModalVisible } = this.state;

    return (
      <dev className="site-layout-content">
        <dev className="ver">
          <Space>
            <Tooltip title={this.state.connectTxt}>
              <ApiFilled style={{ color: this.state.connectColor }} />
            </Tooltip>
            <Tag color="geekblue">v0.2.2</Tag>
            <Tag
              icon={<ReadOutlined />}
              color="#108ee9"
              onClick={this.clickDocTag}
            >
              Document
            </Tag>
            <Radio.Group
              size="small"
              value={locale}
              onChange={this.changeLocale}
            >
              <Radio.Button key="en" value={enUS}>
                English
              </Radio.Button>
              <Radio.Button key="cn" value={zhCN}>
                中文
              </Radio.Button>
            </Radio.Group>

            <Image
              preview={false}
              src="https://img.shields.io/github/stars/pojol/gobot?style=social"
              onClick={this.clickGithubTag}
            />
          </Space>
        </dev>
        <Tabs
          defaultActiveKey="Edit"
          activeKey={this.state.tab}
          onChange={this.changeTab}
        >
          <TabPane tab={lanMap["app.tab.edit"][moment.locale()]} key="Edit">
            <EditPage />
            <TreeModel />
          </TabPane>
          <TabPane tab={lanMap["app.tab.home"][moment.locale()]} key="Home">
            <HomePage />
          </TabPane>
          <TabPane
            tab={lanMap["app.tab.running"][moment.locale()]}
            key="Running"
          >
            <RunningPage />
          </TabPane>
          <TabPane tab={lanMap["app.tab.report"][moment.locale()]} key="Report">
            <Layout>
              <ReportPage />
            </Layout>
          </TabPane>
          <TabPane tab={lanMap["app.tab.prefab"][moment.locale()]} key="Prefab">
            <BotPrefab />
          </TabPane>
          <TabPane tab={lanMap["app.tab.config"][moment.locale()]} key="Config">
            <ConfigPage />
          </TabPane>

        </Tabs>

        <Modal
          visible={isModalVisible}
          onOk={this.modalHandleOk}
          onCancel={this.modalHandleCancel}
        >
          <Tag>e.g. http://123.60.17.61:8888 (Sample driver server address</Tag>
          <Input
            placeholder={lanMap["app.main.modal.input"][moment.locale()]}
            onChange={this.modalConfigChange}
          />
        </Modal>
      </dev>
    );
  }
}
