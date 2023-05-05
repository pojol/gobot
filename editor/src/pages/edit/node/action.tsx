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
  nodeUpdate,
} from "@/models/tree";
import { getDefaultNodeNotifyInfo } from '@/models/node';
import { find } from '@/models/newtree';
import { delay } from '@/utils/timer';

/// <reference path="node.d.ts" />

const { Search } = Input;

export default function ActionTab() {
  const [state, setState] = useState({
    nod: getDefaultNodeNotifyInfo(),
    node_ty: "",
    code: "",
    defaultAlias: "",
    hflex: 0.5,
    wflex: 0.4,
  });

  const [editorState, setEditorState] = useState(null);

  const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);
  const { graphFlex, editFlex } = useSelector((state: RootState) => state.resizeSlice)
  const { nodes } = useSelector((state: RootState) => state.treeSlice)
  const dispatch = useDispatch()

  useEffect(() => {

    delay(100).then(()=>{
      let nod = find(nodes, currentClickNode.id)
      console.info("find", currentClickNode.id, nod)
  
      setState({
        ...state,
        nod: nod,
        code: nod.code,
        defaultAlias: nod.alias,
        node_ty: nod.ty,
      });
    })

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

  const redraw = (wflex: number, hflwx: number) => {
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
