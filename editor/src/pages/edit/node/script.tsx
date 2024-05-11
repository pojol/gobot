import React, { useState, useEffect } from 'react';
import { Input, Button, message, Space } from "antd";

import { useDispatch, useSelector } from 'react-redux';
import { RootState } from "@/models/store";

import PubSub from "pubsub-js";
import Topic from "@/constant/topic";

import { getDefaultNodeNotifyInfo } from '@/models/node';
import { UpdateType, nodeUpdate, find, nodeRedraw } from '@/models/tree';
import { delay } from '@/utils/timer';

import CodeMirror from '@uiw/react-codemirror';
import { StreamLanguage } from '@codemirror/language';
import { lua } from '@codemirror/legacy-modes/mode/lua'
import { xcodeLight, xcodeDark } from '@uiw/codemirror-theme-xcode';
import ThemeType from '@/constant/constant';

/// <reference path="node.d.ts" />

const { Search } = Input;

export default function ScriptTab() {

  const [state, setState] = useState({
    nod: getDefaultNodeNotifyInfo(),
    node_ty: "",
    code: "",
    defaultAlias: "",
  });

  const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);
  const { nodes } = useSelector((state: RootState) => state.treeSlice)
  const dispatch = useDispatch()

  const { themeValue } = useSelector((state: RootState) => state.configSlice)

  useEffect(() => {

    delay(100).then(() => {
      let nod = find(nodes, currentClickNode.id)

      setState({
        ...state,
        nod: nod,
        code: nod.code,
        defaultAlias: nod.alias,
        node_ty: nod.ty,
      });
    })

    return () => {
      // 取消订阅
      PubSub.unsubscribe(Topic.WindowResize);
    };
  }, [currentClickNode])


  const getTheme = () => {
    if (themeValue === ThemeType.Dark) {
      return xcodeDark
    } else {
      return xcodeLight
    }
  }


  const applyClick = () => {
    if (state.nod.id === "") {
      message.warning("节点未被选中");
      return;
    }

    let info = getDefaultNodeNotifyInfo()
    info.id = state.nod.id
    info.code = state.code
    info.alias = state.defaultAlias
    dispatch(nodeUpdate({
      info: info,
      type: [UpdateType.UpdateCode, UpdateType.UpdateAlias]
    }))
    dispatch(nodeRedraw())
  };

  const onChange = (value:any, viewUpdate:any) => {
    setState({
      ...state,
      code: value,
    });
  }

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
        readOnly={false}
        theme={getTheme()}
        extensions={[StreamLanguage.define(lua)]}
        onChange={onChange}
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
