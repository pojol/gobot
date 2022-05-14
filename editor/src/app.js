import { Layout, Tabs, Tag, Radio, Modal, Input, Image, Space } from "antd";
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

import enUS from "antd/lib/locale/en_US";
import zhCN from "antd/lib/locale/zh_CN";
import moment from "moment";
import "moment/locale/zh-cn";
import lanMap from "./locales/lan";
import { PostGetBlob } from "./utils/request";
import Api from "./constant/api";

import { ReadOutlined } from "@ant-design/icons";

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
  }

  componentDidMount() {
    PubSub.subscribe(Topic.FileLoad, (topic, info) => {
      this.setState({ tab: "Edit" });
      PubSub.publish(Topic.FileLoadDraw, [info.Tree]);
    });

    let remote = localStorage.remoteAddr;
    if (remote === "") {
      this.setState({ isModalVisible: true });
    } else {
      this.syncTemplateCode();
    }

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
      }
    });
  };

  syncTemplateCode() {
    PostGetBlob(localStorage.remoteAddr, Api.ConfigGet, {}).then((file) => {
      let reader = new FileReader();
      reader.onload = function (ev) {
        localStorage.CodeTemplate = reader.result;

        PubSub.publish(Topic.ConfigUpdate, {
          key: "code",
          val: reader.result,
        });
      };
      reader.readAsText(file.blob);
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
            <Tag color="geekblue">v0.1.9</Tag>
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
          <TabPane tab={lanMap["app.tab.config"][moment.locale()]} key="Config">
            <ConfigPage />
          </TabPane>
        </Tabs>

        <Modal
          visible={isModalVisible}
          onOk={this.modalHandleOk}
          onCancel={this.modalHandleCancel}
        >
          <Input
            placeholder={lanMap["app.main.modal.input"][moment.locale()]}
            onChange={this.modalConfigChange}
          />
        </Modal>
      </dev>
    );
  }
}
