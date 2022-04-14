import { Layout, Tabs, Tag, Radio } from "antd";
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

const { TabPane } = Tabs;
moment.locale('en');

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      tab: "Edit",
      locale: enUS,
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

  changeLocale = e => {
    const localeValue = e.target.value;
    this.setState({ locale: localeValue });
    
    moment.locale(localeValue.locale);
    
    console.info("moment=>",moment.locale())
  };

  render() {
    const { locale } = this.state;

    return (
      <dev className="site-layout-content">
        <dev className="ver">
          <Tag color="#108ee9">v0.1.5</Tag>
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

      </dev>

    );
  }
}
