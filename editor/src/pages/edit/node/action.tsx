import React, { useState } from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";
import { Input, Button, message, Space } from "antd";

import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";

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

  /*
PubSub.subscribe(Topic.EditPanelCodeMetaResize, (topic, flex) => {
  this.setState({ hflex: flex }, () => {
    this.redraw();
  });
});
*/
  /*
        PubSub.subscribe(Topic.EditPanelEditCodeResize, (topic, flex) => {
          this.setState({ wflex: 1 - flex }, () => {
            this.redraw();
          });
        });
        */

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

    /*
    PubSub.publish(Topic.UpdateNodeParm, {
      parm: {
        id: state.nod.id,
        ty: state.node_ty,
        code: state.code,
        alias: state.defaultAlias,
      },
      notify: true,
    });
    */

    var nod = { ...state.nod };
    //nod.alias = state.defaultAlias;

    //PubSub.publish(Topic.UpdateNode, nod);

    message.success("修改成功");
  };

  const onSearch = (value: any) => {
    if (value !== "") {
      //PubSub.publish(Topic.SearchNode, value);
    }
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

  PubSub.subscribe(Topic.NodeEditorClick, (topic: string, dat: any) => {
    var obj = window.tree.get(dat.id);
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
  });

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
          onSearch={onSearch}
        />
        <Button type="dashed">{state.nod.id}</Button>
      </Space>
    </Space>
  );
}
