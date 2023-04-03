import * as React from "react";
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
import Topic from "@/constant/topic";

const Min = 1;
const Max = 60 * 60 * 1000; // 1 hour

const { Search } = Input;

export default class WaitTab extends React.Component {
    state = {
        nod: {id:"", wait:0},
        node_ty: "WaitNode",
        inputValue: 1,
        defaultAlias:""
    };

    componentDidMount() {
        PubSub.subscribe(Topic.NodeEditorClick, (topic: string, dat: NodeNotifyInfo) => {
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


    onChange = (value:any) => {
        this.setState({
            inputValue: value,
        });
    };

    formatter = (value:any) => {
        return `Delay ${value} ms`;
    };

    applyClick = () => {
        if (this.state.nod.id === "") {
            message.warning("节点未被选中");
            return;
        }

        /*
        PubSub.publish(Topic.UpdateNodeParm, {
            parm: {
                id: this.state.nod.id,
                ty: this.state.node_ty,
                wait: this.state.inputValue,
            },
            notify: true,
        });
*/
        var nod = this.state.nod;
        nod.wait = this.state.inputValue;
        this.setState({ nod: nod });
    };

    render() {
        const { inputValue } = this.state;
        const nod = this.state.nod;

        return (
            <div>

                <Space direction="vertical">
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
                    <Search
                        width={200}
                        enterButton={"apply"}
                        value={this.state.defaultAlias}
                        onSearch={this.applyClick}
                    />
                    <Button type="dashed">{nod.id}</Button>

                </Space>{" "}

            </div>
        );
    }
}
