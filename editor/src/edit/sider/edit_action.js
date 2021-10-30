import React from "react";
import PubSub from "pubsub-js";
import Topic from "../../model/topic";
import ReactJson from "react-json-view";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";

import { Input, Divider, Button, AutoComplete, Tabs, message, Tag, Space } from "antd";
import { NodeTy } from "../../model/node_type";

const { TextArea } = Input;
const { Search } = Input;

export default class ActionTab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      nod: {},
      node_ty: NodeTy.Action,
      code: "",
      defaultAlias: "",
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeEditorClick, (topic, dat) => {
      var obj = window.tree.get(dat.id);

      if (obj !== undefined && obj.ty == this.state.node_ty) {
        let target = { ...obj };
        delete target.pos;
        delete target.children;

        console.info(target)
        this.setState({
          nod: target,
          code: target.code,
          defaultAlias: target.alias,
        });
      } else {
        this.setState({
          nod: {},
        });
      }
    });

    PubSub.subscribe(Topic.EditPlaneResize, (topic, h) =>{
      var nh = h - 100
      this.state.editor.setSize("auto",nh.toString())
    })
  }

  applyClick = () => {
    if (this.state.nod.id === "") {
      message.warning("节点未被选中");
      return;
    }

    PubSub.publish(Topic.UpdateNodeParm, {
      parm: {
        id: this.state.nod.id,
        ty: this.state.node_ty,
        code: this.state.code,
        alias: this.state.defaultAlias,
      },
      notify: true,
    });

    var nod = this.state.nod;
    this.setState({ nod: nod });
  };

  onBeforeChange = (editor, data, value) => {
    this.setState({ code: value });
  };

  onChange = (editor, data, value) => {};

  onDidMount = (editor) => {
    editor.setSize("auto", "400px")
    this.setState({editor:editor})
  }

  onChangeAlias = (e) => {
    this.setState({defaultAlias:e.target.value})
  }

  render() {
    const code = this.state.code;
    const nod = this.state.nod;
    const options = {
      mode: "text/x-lua",
      theme: "solarized dark",
      lineNumbers: true,
    };

    return (
      <div>
        <CodeMirror
          value={code}
          options={options}
          onBeforeChange={this.onBeforeChange}
          onChange={this.onChange}
          editorDidMount={this.onDidMount}
        />
        <Space>
          <Tag color="#55acee">{nod.id}</Tag>
          <Input placeholder="set alias" width={200} value={this.state.defaultAlias} onChange={this.onChangeAlias}/>
          <Button onClick={this.applyClick}>Apply</Button>
        </Space>{" "}
      </div>
    );
  }
}
