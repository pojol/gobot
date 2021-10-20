import {
    Table,
    Tag,
    Space,
    Checkbox,
    InputNumber,
    Divider,
    Button,
    Upload,
    message,
    Input,
    Popconfirm,
    Tooltip,
} from "antd";
import * as React from "react";
import {
    MessageOutlined,
    LikeOutlined,
    StarOutlined,
    InboxOutlined,
    CloudDownloadOutlined,
    SearchOutlined,
    VerticalAlignBottomOutlined,
    DeleteOutlined,
    PlayCircleOutlined,
} from "@ant-design/icons";
import Highlighter from "react-highlight-words";
import Sider from "antd/lib/layout/Sider";
import Config from "../model/config";
import PubSub from "pubsub-js";
import Topic from "../model/topic";
import { Post } from "../model/request";
import { formatTimeStr } from "antd/lib/statistic/utils";
import Api from "../model/api";
import { NodeTy, IsScriptNode } from "../model/node_type";


export default class RunningList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            searchText: "",
            searchedColumn: "",
            runs: {},
            columns: [
                {
                    title: "ID",
                    dataIndex: "id",
                    key: "id",
                },
                {
                    title: "Name",
                    dataIndex: "name",
                    key: "name",
                },
                {
                    title: "Current",
                    dataIndex: "cur",
                    key: "cur",
                },
                {
                    title: "Target",
                    dataIndex: "max",
                    key: "max",
                },
                {
                    title: "Errors",
                    dataIndex: "errors",
                    key: "errors",
                },
            ],
            botLst: [],
        };
    }

    componentDidMount() {
        PubSub.subscribe(Topic.RunningUpdate, (topic, info) => {
            this.refreshBotList();
        });

        this.refreshBotList();
    }

    fillBotList(lst) {
        if (lst) {
            var botlist = [];
            for (var i = 0; i < lst.length; i++) {
                botlist.push({
                    id: lst[i].ID,
                    key: lst[i].ID,
                    name: lst[i].Name,
                    cur: lst[i].Cur,
                    max: lst[i].Max,
                    errors: lst[i].Errors,
                });
            }
            this.setState({ botLst: botlist });
        }
    }

    refreshBotList() {
        this.setState({ botLst: [] });

        Post(window.remote, Api.BotList, {}).then((json) => {
            if (json.Code !== 200) {
                message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
            } else {
                this.fillBotList(json.Body.Lst);
            }
        });
    }

    render() {

        return (
            <div>
                <Table columns={this.state.columns} dataSource={this.state.botLst} />
            </div>
        );
    }
}
