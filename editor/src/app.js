import { Menu, Layout, message, Tabs } from "antd";
import * as React from "react";
import { FileTextOutlined, CloudUploadOutlined } from "@ant-design/icons";
import "antd/dist/antd.css";
import "./app.css";
import GraphView from "./edit/graph/graph";
import Blackboard from "./edit/sider/blackboard";
import TreeModel from "./model/tree_model";
import PubSub from "pubsub-js";
import Topic from "./model/topic";
import BotList from "./home/home";
import TestReport from "./drive/report";
import BotConfig from "./config/config";
import { NodeTy,IsScriptNode } from "./model/node_type";
import EditPlane from "./edit/edit";
import RunningList from "./runing/runing";

const { SubMenu } = Menu;
const { Header, Sider, Content } = Layout;
const { TabPane } = Tabs;

function getValueByElement(elem, tag) {
  for (var i = 0; i < elem.childNodes.length; i++) {
    if (elem.childNodes[i].nodeName === "code") {
      return elem.childNodes[i].childNodes[0].nodeValue;
    }
  }
  return undefined;
}

function parseChildren(xmlnode, children) {
  var nod = {};

  nod.id = xmlnode.getElementsByTagName("id")[0].childNodes[0].nodeValue;
  nod.ty = xmlnode.getElementsByTagName("ty")[0].childNodes[0].nodeValue;

  if (nod.ty === NodeTy.Loop) {
    nod.loop = getValueByElement(xmlnode, "loop")
  } else if (nod.ty === NodeTy.Wait) {
    nod.wait = getValueByElement(xmlnode, "wait")
  } else if (IsScriptNode(nod.ty)) {
    nod.code = getValueByElement(xmlnode, "code");
    nod.alias = getValueByElement(xmlnode, "alias");
  }

  nod.pos = {
    x: parseInt(
      xmlnode.getElementsByTagName("pos")[0].getElementsByTagName("x")[0]
        .childNodes[0].nodeValue
    ),
    y: parseInt(
      xmlnode.getElementsByTagName("pos")[0].getElementsByTagName("y")[0]
        .childNodes[0].nodeValue
    ),
  };

  nod.children = [];
  children.push(nod);

  for (var i = 0; i < xmlnode.childNodes.length; i++) {
    if (xmlnode.childNodes[i].nodeName === "children") {
      parseChildren(xmlnode.childNodes[i], nod.children);
    }
  }
}

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
