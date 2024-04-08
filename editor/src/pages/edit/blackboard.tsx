import React, { useState, useEffect } from 'react';
import ReactJson, { ThemeKeys } from "react-json-view";
import { message, Tabs } from "antd";
import { CodeOutlined, FileSearchOutlined } from "@ant-design/icons";
import { useSelector } from 'react-redux'
import Editor from "react-medium-editor";
import { RootState } from '@/models/store';
import "./blackboard.css";

require("medium-editor/dist/css/medium-editor.css");
require("medium-editor/dist/css/themes/default.css");

import ThemeType from '@/constant/constant';

const { TabPane } = Tabs;

export default function Blackboard() {

    const { threadInfo } = useSelector((state: RootState) => state.debugInfoSlice);
    const [runtimeerr, setRuntimeerr] = useState("");
    const [change, setChange] = useState("")
    const [active, setActive] = useState("2")
    const [jsontheme, setJsontheme] = useState<ThemeKeys>('google')
    const { themeValue } = useSelector((state: RootState) => state.configSlice)

    const metainfo = useSelector((state: RootState) => state.debugInfoSlice.metaInfo)

    useEffect(() => {
        let msg = ""
        let haveerr = false

        if (themeValue === ThemeType.Dark) {
            setJsontheme("google")
        } else {
            setJsontheme("rjv-default")
        }

        try {
            threadInfo.forEach(element => {
                msg += "<b>Thread[" + element.number + "]</b>\n"

                if (element.errmsg !== "") {
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
            setActive("3")
            setRuntimeerr(msg)
        } else {
            if (active === "3") {
                setActive("2")
            }

            setRuntimeerr("")
            setChange(msg)
        }

    }, [threadInfo, themeValue])

    const clickTab = (activeKey: string, e: React.KeyboardEvent<Element> | React.MouseEvent<Element, MouseEvent>) => {
        setActive(activeKey)
    }

    return (
        <div className="scroll-patch">
            <Tabs activeKey={active} onTabClick={clickTab}>
                <TabPane
                    tab={
                        <span>
                            <FileSearchOutlined />
                            Blackboard
                        </span>
                    }
                    key="1"
                >
                    <ReactJson
                        name=""
                        src={JSON.parse(metainfo)}
                        theme={jsontheme}
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
                        text={change}
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
                        text={runtimeerr}
                    />
                </TabPane>
            </Tabs>
        </div>
    );

}
