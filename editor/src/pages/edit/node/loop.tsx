import React, { useState } from "react";
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

import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";

const Min = 0;
const Max = 10000;

const { Search } = Input;

export default function LoopTab() {
  const [state, setState] = useState({
    inputValue: 0,
    nod: { id: "", loop: 0 },
    node_ty: "LoopNode",
    defaultAlias: "",
  });

  const onChange = (value: any) => {
    setState({
      ...state,
      inputValue: value,
    });
  };

  const formatter = (value: any) => {
    if (value === 0) {
      return `endless`;
    } else {
      return `loop ${value} times`;
    }
  };

  const applyClick = () => {
    if (state.nod.id === "") {
      message.warning("节点未被选中");
      return;
    }

    /*
        PubSub.publish(Topic.UpdateNodeParm, {
            parm: {
                id: this.state.nod.id,
                ty: this.state.node_ty,
                loop: this.state.inputValue,
            },
            notify: true,
        });
        */

    var nod = state.nod;
    nod.loop = state.inputValue;
    setState({
      ...state,
      nod: nod,
    });
  };

  PubSub.subscribe(Topic.NodeEditorClick, (topic: string, dat: any) => {
    console.info("Topic loop, componentDidMount");
    var obj = window.tree.get(dat.id);
    if (obj !== undefined && obj.ty === state.node_ty) {
      let target = { ...obj };
      delete target.pos;
      delete target.children;

      setState({
        ...state,
        nod: target,
        inputValue: target.loop,
      });
    } else {
      setState({
        ...state,
        nod: { id: "", loop: 0 },
        inputValue: 0,
      });
    }
  });

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
              value={
                typeof state.inputValue === "number" ? state.inputValue : 0
              }
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
          enterButton={"Apply"}
          value={state.defaultAlias}
          onSearch={applyClick}
        />
        <Button type="dashed">{state.nod.id}</Button>
      </Space>{" "}
    </div>
  );
}
