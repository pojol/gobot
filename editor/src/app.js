import {  Layout, Tabs } from "antd";
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

const { TabPane } = Tabs;



export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      tab: "Edit",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.FileLoad, (topic, info) => {
      this.setState({ tab: "Edit" });
      PubSub.publish(Topic.FileLoadGraph, info.Tree);
    });
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

  render() {
    return (
      <Layout>
        <Layout>
          <Tabs
            defaultActiveKey="Edit"
            activeKey={this.state.tab}
            onChange={this.changeTab}
          >
            <TabPane tab="Edit" key="Edit">
              <EditPlane />
              <TreeModel />
            </TabPane>
            <TabPane tab="Home" key="Home">
              <BotList />
            </TabPane>
            <TabPane tab="Running" key="Running">
              <RunningList/>
            </TabPane>
            <TabPane tab="Report" key="Report">
              <Layout>
                <TestReport />
              </Layout>
            </TabPane>
            <TabPane tab="Config" key="Config">
              <BotConfig />
            </TabPane>
          </Tabs>
        </Layout>
      </Layout>
    );
  }
}
