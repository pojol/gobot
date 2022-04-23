import { Layout, Tabs, Tag, Radio, Modal, Input } from "antd";
import * as React from "react";
import "antd/dist/antd.css";
import "./app.css";
import TreeModel from "./model/tree_model";
import PubSub from "pubsub-js";
import Topic from "./model/topic";
import BotList from "./home/home";
import TestReport from "./drive/report";
import BotConfig from "./config/config";
import EditPlane from "./edit/edit";
import RunningList from "./runing/runing";

import enUS from 'antd/lib/locale/en_US';
import zhCN from 'antd/lib/locale/zh_CN';
import moment from 'moment';
import 'moment/locale/zh-cn';
import lanMap from "./config/lan";
import { PostGetBlob } from "./model/request";
import Api from "./model/api";

const { TabPane } = Tabs;
moment.locale('en');

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
        resizeTimer = null
        PubSub.publish(Topic.WindowResize, {});
      }, 100)
    }
  }

  componentDidMount() {
    PubSub.subscribe(Topic.FileLoad, (topic, info) => {
      this.setState({ tab: "Edit" });
      PubSub.publish(Topic.FileLoadDraw, [info.Tree]);
    });

    let remote = localStorage.remoteAddr
    if (remote === "") {
      this.setState({isModalVisible: true})
    } else {
      this.syncTemplateCode()
    }

    window.addEventListener('resize', this.resizeHandler, false)
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

    PostGetBlob(localStorage.remoteAddr, Api.ConfigGet, {}).then(
      (file) => {
        let reader = new FileReader();
        reader.onload = function(ev) {
          localStorage.CodeTemplate = reader.result

          PubSub.publish(Topic.ConfigUpdate, {
            key : "code",
            val : reader.result,
          })
        }
        reader.readAsText(file.blob);
      }
    )

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
      localStorage.remoteAddr = this.state.modalConfig

      this.syncTemplateCode()
    }
  };

  modalHandleCancel = () => {
    this.setState({ isModalVisible: false });
  };


  changeLocale = e => {
    const localeValue = e.target.value;
    this.setState({ locale: localeValue });
    
    moment.locale(localeValue.locale);
    
    console.info("moment=>",moment.locale())
  };

  render() {
    const { locale, isModalVisible } = this.state;

    return (
      <dev className="site-layout-content">
        <dev className="ver">
          <Tag color="#108ee9">v0.1.6</Tag>
          <Radio.Group value={locale} onChange={this.changeLocale}>
            <Radio.Button key="en" value={enUS}>
              English
            </Radio.Button>
            <Radio.Button key="cn" value={zhCN}>
              中文
            </Radio.Button>
          </Radio.Group>
        </dev>
        <Tabs
          defaultActiveKey="Edit"
          activeKey={this.state.tab}
          onChange={this.changeTab}
        >
          <TabPane tab={lanMap["app.tab.edit"][moment.locale()]} key="Edit">
            <EditPlane />
            <TreeModel />
          </TabPane>
          <TabPane tab={lanMap["app.tab.home"][moment.locale()]} key="Home">
            <BotList />
          </TabPane>
          <TabPane tab={lanMap["app.tab.running"][moment.locale()]} key="Running">
            <RunningList />
          </TabPane>
          <TabPane tab={lanMap["app.tab.report"][moment.locale()]} key="Report">
            <Layout>
              <TestReport />
            </Layout>
          </TabPane>
          <TabPane tab={lanMap["app.tab.config"][moment.locale()]} key="Config">
            <BotConfig />
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
