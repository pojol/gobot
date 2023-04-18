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
import { Root } from 'react-dom/client';

/// <reference path="node.d.ts" />

const { Search } = Input;

export default function ActionTab() {
  const [state, setState] = useState({
    nod: { id: "" },
    node_ty: "",
    code: "",
    defaultAlias: "",
    hflex: 0,
    wflex: 0,
    editor: undefined,
  });

  const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);
  const { graphFlex, editFlex } = useSelector((state:RootState)=> state.resizeSlice)
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
  },[currentClickNode])

  useEffect(()=>{
    setState({
      ...state,
      wflex: 1 - graphFlex,
    });
    
    // redraw

  }, [graphFlex])

  useEffect(()=>{
    setState({
      ...state,
      hflex: editFlex,
    });

    // redraw

  }, [editFlex])

  /*
    PubSub.subscribe(Topic.WindowResize, () => {
      this.redraw();
    });
    */

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
          setState({
            ...state,
            wflex: 0.4,
            hflex: 0.5,
            editor: editor,
          });

          var width = document.documentElement.clientWidth * 0.4 - 18;
          var height = document.documentElement.clientHeight * 0.5 - 38;

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
