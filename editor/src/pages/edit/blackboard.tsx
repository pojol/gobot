import React, { useState, useEffect } from 'react';
import ReactJson, { ThemeKeys } from "react-json-view";
import { message, Tabs } from "antd";
import { CodeOutlined, FileSearchOutlined } from "@ant-design/icons";
import { useSelector } from 'react-redux'
import Editor from "react-medium-editor";
import { RootState } from '@/models/store';
import "./blackboard.css";
const { Post } = require("../../utils/request");
import Api from "../../constant/api";

require("medium-editor/dist/css/medium-editor.css");
require("medium-editor/dist/css/themes/default.css");
import ThemeType from '@/constant/constant';
import { TaskTimer } from 'tasktimer';

const { TabPane } = Tabs;

const stdoutstr = async (botid: string) => {
    let info = "";
    try {
        const json = await Post(localStorage.remoteAddr, Api.RuntimeInfo, { ID: botid });
        if (json.Code == 200) {
            info = json.Body.Msg;
        }
    } catch (err) {
        console.error("Error:", err);
    }

    if (info !== "") {
        console.info("post runtime.info", info);
    }

    return info;
};


export function Stdout() {

    const [change, setChange] = useState(String.raw`
                                __              __      
                               /\ \            /\ \__   
                       __     ___\ \ \____    ___\ \ ,_\  
                     /'_ '\  / __'\ \ '__'\  / __'\ \ \  
                    /\ \L\ \/\ \L\ \ \ \L\ \/\ \L\ \ \ \_ 
                    \ \____ \ \____/\ \_,__/\ \____/\ \__\
                     \/___L\ \/___/  \/___/  \/___/  \/__/
                       /\____/                            
                       \_/__/           <b>v0.4.4</b>                 
    `)
    const { themeValue } = useSelector((state: RootState) => state.configSlice)
    const { currentDebugBot } = useSelector((state: RootState) => state.treeSlice)

    useEffect(() => {

        const newTimer = new TaskTimer(500);
        newTimer.on('tick', async () => {
            if (currentDebugBot === "") {
                return
            }

            let newmsg = await stdoutstr(currentDebugBot);
            if (newmsg !== "") {
                setChange(prev => {
                    let oldmsg = newmsg + "\n" + prev;
                    return oldmsg
                })
            }
        });
        newTimer.start();

        return () => {
            // 清理定时器
            if (newTimer) {
                newTimer.stop();
            }
        }

    }, [themeValue, currentDebugBot])

    return (
        <div>
            <Tabs activeKey={"1"}>

                <TabPane
                    tab={
                        <span>
                            <CodeOutlined />
                            Stdout
                        </span>
                    }
                    key="1"
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

            </Tabs>
        </div>
    );

}

export function Blackboard() {
    const [jsontheme, setJsontheme] = useState<ThemeKeys>('google')
    const { themeValue } = useSelector((state: RootState) => state.configSlice)
    const metainfo = useSelector((state: RootState) => state.debugInfoSlice.metaInfo)

    useEffect(() => {

        if (themeValue === ThemeType.Dark) {
            setJsontheme("google")
        } else {
            setJsontheme("rjv-default")
        }

    }, [themeValue])

    useEffect(() => {

    }, [metainfo])

    return (
        <div className='ant-tabs-tabpane'>
            <Tabs activeKey={"1"}>
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
                        name="bot"
                        src={JSON.parse(metainfo)}
                        theme={jsontheme}
                        enableClipboard={false}
                        displayDataTypes={false}
                    ></ReactJson>
                </TabPane>


            </Tabs>
        </div>
    )
}
