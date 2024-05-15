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

export function Stdout() {

    const { threadInfo } = useSelector((state: RootState) => state.debugInfoSlice);
    const [runtimeerr, setRuntimeerr] = useState("");
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

    useEffect(() => {
        let msg = ""
        let haveerr = false

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

            if (msg !== "") {
                setChange(msg)
            }
        } catch (err) {
        }

    }, [threadInfo, themeValue])

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
