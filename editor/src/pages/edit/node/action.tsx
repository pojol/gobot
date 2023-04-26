import React, { useState, useEffect } from 'react';
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";
import { Input, Button, message, Space } from "antd";

import { useDispatch, useSelector } from 'react-redux';
import { RootState } from "@/models/store";

import PubSub from "pubsub-js";
import Topic from "@/constant/topic";

import {
  getDefaultNodeNotifyInfo,
  nodeUpdate,
} from "@/models/tree";

/// <reference path="node.d.ts" />

const { Search } = Input;

export default function ActionTab() {
  const [state, setState] = useState({
    nod: { id: "" },
    node_ty: "",
    code: "",
    defaultAlias: "",
    hflex: 0.5,
    wflex: 0.4,
  });

  const [editorState, setEditorState] = useState(null);

  const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);
  const { graphFlex, editFlex } = useSelector((state: RootState) => state.resizeSlice)
  const dispatch = useDispatch()

  useEffect(() => {
    var obj = window.tree.get(currentClickNode.id);
    if (obj !== undefined) {
      let target = { ...obj };
      delete target.pos;
      delete target.children;

      setState({
        ...state,
        nod: target,
        code: target.code,
        defaultAlias: target.alias,
        node_ty: target.ty,
      });
    } else {
      setState({
        ...state,
        nod: { id: "" },
      });
    }

    PubSub.subscribe(Topic.WindowResize, () => {
        redraw(state.wflex, state.hflex)
    });

    return () => {
      // 取消订阅
      PubSub.unsubscribe(Topic.WindowResize);
  };
  }, [currentClickNode])

  useEffect(() => {
    setState({
      ...state,
      wflex: 1 - graphFlex,
    });

    console.info("wflex", graphFlex)
    redraw((1 - graphFlex), state.hflex)
  }, [graphFlex])

  useEffect(() => {
    setState({
      ...state,
      hflex: editFlex,
    });

    console.info("hflex", editFlex)
    redraw(state.wflex, editFlex)
  }, [editFlex])

  const redraw = (wflex : number, hflwx : number) => {
    if (editorState !== null) {
      var width = document.documentElement.clientWidth * wflex - 18;
      var height = document.documentElement.clientHeight * hflwx - 40;
      console.info("redraw", wflex, hflwx)
      editorState.setSize(width.toString() + "px", height.toString() + "px");
    }
  }

  const applyClick = () => {
    if (state.nod.id === "") {
      message.warning("节点未被选中");
      return;
    }

    let info = getDefaultNodeNotifyInfo()
    info.id = state.nod.id
    info.ty = state.node_ty
    info.code = state.code
    info.alias = state.defaultAlias
    info.notify = true
    dispatch(nodeUpdate(info))
    PubSub.publish(Topic.UpdateNodeParm, info)
  };

  const handleChange = (editor: any, data: any, value: any) => {
    setState({
      ...state,
      code: value,
    });
  };

  const handleAliasChange = (event: any) => {
    setState({
      ...state,
      defaultAlias: event.target.value,
    });
  };

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      <CodeMirror
        value={state.code}
        onBeforeChange={(editor, data, value) =>
          handleChange(editor, data, value)
        }
        options={{
          mode: "lua",
          theme: localStorage.codeboxTheme,
          lineNumbers: true,
        }}
        editorDidMount={(editor) => {
          setEditorState(editor);

          var width = document.documentElement.clientWidth * state.wflex - 18;
          var height = document.documentElement.clientHeight * state.hflex - 40;

          editor.setSize(width.toString() + "px", height.toString() + "px");
        }}
      />
      <Space>
        <Search
          width={200}
          value={state.defaultAlias}
          enterButton={"Apply"}
          onChange={handleAliasChange}
          onSearch={applyClick}
        />
        <Button type="dashed">{state.nod.id}</Button>
      </Space>
    </Space>
  );
}
