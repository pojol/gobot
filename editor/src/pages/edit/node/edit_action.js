import React from "react";
import PubSub from "pubsub-js";
import Topic from "../../../constant/topic";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/mode/lua/lua";
import { Input, Button, message, Space } from "antd";
import {formatText} from 'lua-fmt';

import moment from "moment";
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
      hflex: 0,
      wflex: 0,
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeEditorClick, (topic, dat) => {
      var obj = window.tree.get(dat.id);
      console.info("click", dat.id, obj)
      if (obj !== undefined) {
        let target = { ...obj };
        delete target.pos;
        delete target.children;

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

    PubSub.subscribe(Topic.EditPanelCodeMetaResize, (topic, flex) => {
      this.setState({ hflex: flex }, () => {
        this.redraw();
      });
    });

    PubSub.subscribe(Topic.EditPanelEditCodeResize, (topic, flex) => {
      this.setState({ wflex: 1 - flex }, () => {
        this.redraw();
      });
    });

    PubSub.subscribe(Topic.WindowResize, () => {
      this.redraw();
    });
  }

  redraw() {
    var width = document.body.clientWidth * this.state.wflex - 2;
    var height = document.body.clientHeight * this.state.hflex - 36;

    this.state.editor.setSize(
      width.toString() + "px",
      height.toString() + "px"
    );
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
    this.setState({ editor: editor, wflex: 0.4, hflex: 0.5 }, () => {
      var width, height;
      var dimensions = this.props.dimensions;

      if (dimensions.width !== "100%") {
        width = dimensions.width - 2;
        height = dimensions.height - 38;

        this.setState({
          wflex: dimensions.width / document.body.clientWidth,
          hflex: dimensions.height / document.body.clientHeight,
        });
      } else {
        width = document.body.clientWidth * this.state.wflex - 2;
        height = document.body.clientHeight * this.state.hflex - 38;
      }

      console.info("action init", dimensions, "w", width, "h", height);
      this.state.editor.setSize(
        width.toString() + "px",
        height.toString() + "px"
      );
    });
  };

  onChangeAlias = (e) => {
    this.setState({ defaultAlias: e.target.value });
  };

  clickFmtBtn = (e) => {
    let old = this.state.code;
    this.setState({code:formatText(old)})
  };

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
          <Button onClick={this.clickFmtBtn}>fmt</Button>
          <Button type="dashed">{nod.id}</Button>
        </Space>{" "}
      </div>
    );
  }
}
