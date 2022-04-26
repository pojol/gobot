import React from "react";
import {
  Slider,
  InputNumber,
  Row,
  Col,
  Button,
  message,
  Tag,
  Space,
  Input,
} from "antd";
import PubSub from "pubsub-js";
import Topic from "../../model/topic";


import moment from 'moment';
import lanMap from "../../config/lan";

const Min = 1;
const Max = 60 * 60 * 1000; // 1 hour

const { Search } = Input;


export default class WaitTab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      nod: {},
      node_ty: "WaitNode",
      inputValue: 1,
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.NodeEditorClick, (topic, dat) => {
      var obj = window.tree.get(dat.id);

      if (obj !== undefined && obj.ty === this.state.node_ty) {
        let target = { ...obj };
        delete target.pos;
        delete target.children;

        this.setState({
          nod: target,
        });
      } else {
        this.setState({
          nod: {},
          inputValue: 1,
        });
      }
    });
  }

  onChange = (value) => {
    this.setState({
      inputValue: value,
    });
  };

  formatter = (value) => {
    return `Delay ${value} ms`;
  };

  applyClick = () => {
    if (this.state.node_id === "") {
      message.warning("节点未被选中");
      return;
    }

    PubSub.publish(Topic.UpdateNodeParm, {
      parm: {
        id: this.state.nod.id,
        ty: this.state.node_ty,
        wait: this.state.inputValue,
      },
      notify: true,
    });

    var nod = this.state.nod;
    nod.wait = this.state.inputValue;
    this.setState({ nod: nod });
  };

  render() {
    const { inputValue } = this.state;
    const nod = this.state.nod;

    return (
      <div>
        <Row>
          <Col span={12}>
            <Slider
              tipFormatter={this.formatter}
              min={Min}
              max={Max}
              onChange={this.onChange}
              value={typeof inputValue === "number" ? inputValue : 1}
            />
          </Col>
          <Col span={4}>
            <InputNumber
              min={Min}
              max={Max}
              style={{ margin: "0 26px" }}
              value={inputValue}
              onChange={this.onChange}
            />
          </Col>
        </Row>
        <Space>
          <Button type="dashed">{nod.id}</Button>
          <Search
            placeholder={lanMap["app.edit.tab.placeholder"][moment.locale()]}
            width={200}
            enterButton={lanMap["app.edit.tab.apply"][moment.locale()]}
            value={this.state.defaultAlias}
            onChange={this.onChangeAlias}
            onSearch={this.applyClick}
          />
        </Space>{" "}
      </div>
    );
  }
}
