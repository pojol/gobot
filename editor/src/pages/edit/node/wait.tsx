import React, { useState, useEffect } from "react";
import {
    InputNumber,
    Row,
    Col,
    Button,
    message,
    Slider,
    Space,
    Input,
} from "antd";
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from "@/models/store";

import PubSub from "pubsub-js";
import Topic from "@/constant/topic";

import { getDefaultNodeNotifyInfo } from "@/models/node";
import { delay } from "@/utils/timer";
import { UpdateType, find, nodeUpdate } from "@/models/newtree";

const Min = 1;
const Max = 60 * 60 * 1000; // 1 hour

const { Search } = Input;

export default function WaitTab() {

    const [state, setState] = useState({
        nod: getDefaultNodeNotifyInfo(),
        node_ty: "WaitNode",
        inputValue: 1,
        defaultAlias: ""
    });
    const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);
    const { nodes } = useSelector((state: RootState) => state.treeSlice)

    const dispatch = useDispatch()

    useEffect(() => {

        delay(100).then(() => {

            let nod = find(nodes, currentClickNode.id)
            setState({
                ...state,
                nod: nod,
                inputValue: nod.wait
            })
        })

    }, [currentClickNode])


    const onChange = (value: any) => {
        setState({
            ...state,
            inputValue: value,
        })
    };

    const formatter = (value: any) => {
        return `Delay ${value} ms`;
    };

    const applyClick = () => {
        if (state.nod.id === "") {
            message.warning("节点未被选中");
            return;
        }

        let info = getDefaultNodeNotifyInfo()
        info.id = state.nod.id
        info.wait = state.inputValue
        dispatch(nodeUpdate({
            info: info,
            type: [UpdateType.UpdateWait]
        }))
    };

    return (
        <div>

            <Space direction="vertical">
                <Row>
                    <Col span={12}>
                        <Slider
                            tipFormatter={formatter}
                            min={Min}
                            max={Max}
                            onChange={onChange}
                            value={typeof state.inputValue === "number" ? state.inputValue : 1}
                        />
                    </Col>
                    <Col span={4}>
                        <InputNumber
                            min={Min}
                            max={Max}
                            style={{ margin: "0 26px" }}
                            value={state.inputValue}
                            onChange={onChange}
                        />
                    </Col>
                </Row>
                <Search
                    width={200}
                    enterButton={"apply"}
                    value={state.defaultAlias}
                    onSearch={applyClick}
                />
                <Button type="dashed">{state.nod.id}</Button>

            </Space>{" "}

        </div>
    );
}
