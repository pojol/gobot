import React from "react";
import ReactJson from "react-json-view";
import { message, Tabs } from "antd";
import { CodeOutlined, FileSearchOutlined } from "@ant-design/icons";

import Editor from "react-medium-editor";

import "./blackboard.css";

require("medium-editor/dist/css/medium-editor.css");
require("medium-editor/dist/css/themes/default.css");

const { TabPane } = Tabs;
import PubSub from "pubsub-js";
import Topic from "@/constant/topic";

export default class Blackboard extends React.Component {

    state = {
        metadata: {},
        change: "",
        runtimeerr: "",
        active: "2",
    };

    componentDidMount() {
        /*
        PubSub.subscribe(Topic.UpdateBlackboard, (topic, info) => {
            try {
                var blackboard = JSON.parse(info);
                this.setState({ metadata: blackboard });
            } catch (err) {
                message.warning("blackboard parse info err");
            }
        });
        */

        /*
        PubSub.subscribe(Topic.UpdateChange, (topic, threadInfo) => {
            let msg = ""
            let haveerr = false

            try {
                threadInfo.forEach(element => {
                    msg += "<b>Thread[" + element.number + "]</b>\n"

                    if (element.errmsg !== "") {
                        PubSub.publish(Topic.Focus, [element.errmsg.substr(0, 36)])
                        msg += element.errmsg
                        msg += "------------------------------\n"
                        haveerr = true
                        throw new Error();
                    }

                    let changemsg = "{}"
                    if (element.change !== "") {
                        changemsg = element.change
                    }

                    try {
                        msg += JSON.stringify(JSON.parse(changemsg), null, 2) + "\n"
                    } catch (error) {
                        console.warn(error)
                        msg += changemsg + "\n"
                    }

                    msg += "------------------------------\n"
                })
            } catch (err) {
            }

            if (haveerr) {
                this.setState({ runtimeerr: msg, active: "3" });
            } else {
                this.setState({ change: msg });
            }

        });
        */

        /*
        PubSub.subscribe(Topic.Upload, (topic, info) => {
            this.setState({ change: "" });
        });
        */

        /*
        PubSub.subscribe(Topic.Upload, (topic, info) => {
            this.setState({ metadata: JSON.parse("{}") });
        });
        */

        PubSub.subscribe(Topic.DebugCreate, (topic: string, info: any) => {
            this.setState({ metadata: JSON.parse("{}"), runtimeerr: "" });
            console.info(this.state.runtimeerr)
        });

    }

    clickTab = (activeKey: string, e: React.KeyboardEvent<Element> | React.MouseEvent<Element, MouseEvent>) => {
        this.setState({ active: e })
    }

    render() {
        return (
            <div className="scroll-patch">
                <Tabs activeKey={this.state.active} onTabClick={this.clickTab}>
                    <TabPane
                        tab={
                            <span>
                                <FileSearchOutlined />
                                Meta
                            </span>
                        }
                        key="1"
                    >
                        <ReactJson
                            name=""
                            src={this.state.metadata}
                            theme={"rjv-default"}
                            enableClipboard={false}
                            displayDataTypes={false}
                        ></ReactJson>
                    </TabPane>
                    <TabPane
                        tab={
                            <span>
                                <CodeOutlined />
                                Response
                            </span>
                        }
                        key="2"
                    >
                        <Editor
                            tag="pre"
                            //https://github.com/yabwe/medium-editor/blob/d113a74437fda6f1cbd5f146b0f2c46288b118ea/OPTIONS.md#disableediting
                            options={{
                                placeholder: { text: "", hideOnClick: true },
                                disableEditing: true,
                            }}
                            text={this.state.change}
                        />
                    </TabPane>
                    <TabPane
                        tab={
                            <span>
                                <CodeOutlined />
                                RuntimeError
                            </span>
                        }
                        key="3"
                    >
                        <Editor
                            tag="pre"
                            //https://github.com/yabwe/medium-editor/blob/d113a74437fda6f1cbd5f146b0f2c46288b118ea/OPTIONS.md#disableediting
                            options={{
                                placeholder: { text: "", hideOnClick: true },
                                disableEditing: true,
                            }}
                            text={this.state.runtimeerr}
                        />
                    </TabPane>
                </Tabs>
            </div>
        );
    }
}
