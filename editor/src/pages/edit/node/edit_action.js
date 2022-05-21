import React from "react";
import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";

import { Input, Button, message, Space } from "antd";

import moment from 'moment';
import lanMap from "../../../locales/lan";

const { Search } = Input;


export default class ActionTab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      nod: {},
      node_ty: "",
      code: "",
      defaultAlias: "",
      wratio: 0.4,
      hratio: 0.53
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeEditorClick, (topic, dat) => {
      var obj = window.tree.get(dat.id);

      if (obj !== undefined) {
        let target = { ...obj };
        delete target.pos;
        delete target.children;

        console.info("click", target)
        this.setState({
          nod: target,
          code: target.code,
          defaultAlias: target.alias,
          node_ty: target.ty,
        });
      } else {
        this.setState({
          nod: {},
        });
      }
    });

    PubSub.subscribe(Topic.EditPlaneCodeMetaResize, (topic, h) => {

      let hratio = 1 - ((document.body.clientHeight - h + 88) / document.body.clientHeight)
      console.info(hratio)

      this.setState({ hratio: hratio })
      this.state.editor.setSize((document.body.clientWidth * this.state.wratio).toString() + "px", (document.body.clientHeight * this.state.hratio).toString() + "px")
    })

    PubSub.subscribe(Topic.EditPlaneEditCodeResize, (topic, w) => {
      let wratio = ((document.body.clientWidth - w) / document.body.clientWidth)
      //console.info(wratio, (document.body.clientWidth * this.state.wratio).toString())

      this.setState({ wratio: wratio })
      this.state.editor.setSize((document.body.clientWidth * this.state.wratio).toString() + "px", (document.body.clientHeight * this.state.hratio).toString() + "px")
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

  onDidMount = (editor) => {
    editor.setSize((document.body.clientWidth * this.state.wratio).toString() + "px", (document.body.clientHeight * this.state.hratio).toString() + "px")
    console.info("document.body.clientWidth", document.body.clientWidth, document.body.clientWidth * this.state.wratio)
    this.setState({ editor: editor })
  }

  onChangeAlias = (e) => {
    this.setState({ defaultAlias: e.target.value })
  }

  render() {
    const code = this.state.code;
    const nod = this.state.nod;
    const options = {
      mode: "text/x-lua",
      theme: localStorage.theme,
      lineNumbers: true,
    };

    return (
      <div>
        <CodeMirror
          value={code}
          options={options}
          onBeforeChange={this.onBeforeChange}
          editorDidMount={this.onDidMount}
        />
        <Space>
          <Search
            placeholder={lanMap["app.edit.tab.placeholder"][moment.locale()]}
            width={200}
            enterButton={lanMap["app.edit.tab.apply"][moment.locale()]}
            value={this.state.defaultAlias}
            onChange={this.onChangeAlias}
            onSearch={this.applyClick}
          />
          <Button type="dashed">{nod.id}</Button>
        </Space>{" "}
      </div>
    );
  }
}
