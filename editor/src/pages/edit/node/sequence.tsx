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

/// <reference path="node.d.ts" />

import PubSub from "pubsub-js";
import Topic from "@/constant/topic";

export default class SequenceTab extends React.Component {
    state = {
        node_id: ""
    };

    componentDidMount() {
        PubSub.subscribe(Topic.NodeEditorClick, (topic: string, dat: NodeNotifyInfo) => {
            this.setState({ node_id: dat.id })
        })
    }

    render() {
        return (
            <div>
                <Button type="dashed">{this.state.node_id}</Button>
            </div>
        );
    }
}
