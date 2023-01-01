import React, { useRef, useState } from 'react';
import {
    Space, Table, Button, message, Input, Popconfirm,
    Tooltip, Modal
} from 'antd';
import PubSub from "pubsub-js";
import { formatText } from 'lua-fmt';
import {
    FilterTwoTone
} from "@ant-design/icons";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/theme/abcdef.css";
import "codemirror/theme/ayu-dark.css";
import "codemirror/theme/yonce.css";
import "codemirror/theme/neo.css";
import "codemirror/theme/zenburn.css";
import "codemirror/mode/lua/lua";

import Topic from "../../constant/topic";
import Api from "../../constant/api";
import moment from 'moment';
import lanMap from "../../locales/lan";

import HomeTagGroup from '../home/home_tags';

import { Post, PostGetBlob, PostBlob } from "../../utils/request";

import "antd/dist/antd.css";
// You will need to import the styles separately
// You probably want to do this just once during the bootstrapping phase of your application.
import "react-reflex/styles.css";

// then you can import the components
import { ReflexContainer, ReflexSplitter, ReflexElement } from "react-reflex";



export default class BotPrefab extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: [],
            code: "",
            selectedKey: "",
            isModalVisible: false,
            newPrefabName: "",
            searchText: "",
            setSearchText: "",
            searchedColumn: "",
            searchInput: "",
        };
    }

    componentDidMount() {
        this.syncConfig()
    }

    handleSearch = (selectedKeys, confirm, dataIndex) => {
        confirm();
        this.setState({ setSearchText: selectedKeys[0] })
    };

    getColumnSearchProps = (dataIndex) => ({
        filterDropdown: ({ setSelectedKeys, selectedKeys, confirm }) => (
            <div
                style={{
                    padding: 8,
                }}
            >
                <Input
                    placeholder={`Search ${dataIndex}`}
                    value={selectedKeys[0]}
                    onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
                    onPressEnter={() => this.handleSearch(selectedKeys, confirm, dataIndex)}
                    style={{
                        marginBottom: 8,
                        display: 'block',
                    }}
                />
            </div>
        ),
        filterIcon: (filtered) => (
            <FilterTwoTone />
        ),
        onFilter: (value, record) =>
            record[dataIndex].toString().toLowerCase().includes(value.toLowerCase()),
        render: (text) => <a>{text.toLowerCase()}</a>
    })

    syncConfig = () => {
        Post(localStorage.remoteAddr, Api.ConfigList, {}).then((json) => {
            if (json.Code !== 200) {
                message.error(
                    "get config list fail:" + String(json.Code) + " msg: " + json.Msg
                );
            } else {
                let lst = json.Body.Lst;
                var counter = 0;
                let dat = []

                let callback = () => {
                    dat.sort(function (a, b) {
                        if (a.name < b.name) { return -1; }
                        if (a.name > b.name) { return 1; }
                        return 0;
                    });

                    this.setState({ data: dat })
                };

                lst.forEach(function (element) {

                    PostGetBlob(localStorage.remoteAddr, Api.ConfigGet, element).then(
                        (file) => {
                            let reader = new FileReader();
                            reader.onload = function (ev) {
                                let lowElement = element.toLowerCase()
                                console.info("element", lowElement)
                                if (lowElement !== "system" && lowElement != "global") {
                                    var jobj = JSON.parse(reader.result);
                                    dat.push({ key: jobj["title"].toLowerCase(), name: jobj["title"].toLowerCase(), code: jobj["content"] })
                                }

                                counter++;
                                if (counter === lst.length) {
                                    callback()
                                    PubSub.publish(Topic.ConfigUpdateAll, {})
                                }
                            };

                            reader.readAsText(file.blob);
                        }
                    );

                });
            }
        });
    }

    onBeforeChange = (editor, data, value) => {
        this.setState({ code: value });
    };

    uploadPrefab = (name, code) => {
        var templatecode = JSON.stringify({
            title: name,
            content: code,
            key: name,
        });
        var blob = new Blob([templatecode], {
            type: "application/json",
        });

        console.info(name, "apply config code", templatecode);

        PostBlob(
            localStorage.remoteAddr,
            Api.ConfigUpload,
            name,
            blob
        ).then((json) => {
            if (json.Code !== 200) {
                message.error(
                    "upload fail:" + String(json.Code) + " msg: " + json.Msg
                );
            } else {
                message.success(name + " upload succ");
                window.config.set(name, templatecode);
            }

            this.syncConfig()
        });
    }

    onHandleAdd = () => {
        this.setState({ isModalVisible: true })
    }

    onHandleUpdate = () => {
        let targetKey = this.state.selectedKey
        if (targetKey === "") {
            return
        }

        this.uploadPrefab(targetKey, this.state.code)
    }

    onHandleRemove = () => {
        let targetKey = this.state.selectedKey
        if (targetKey === "") {
            return
        }

        Post(localStorage.remoteAddr, Api.ConfigRemove, {
            Name: targetKey,
        }).then((json) => {
            if (json.Code !== 200) {
                message.error(
                    "remove prefab err:" + String(json.Code) + " msg: " + json.Msg
                );
            } else {
                window.config.delete(targetKey);
                PubSub.publish(Topic.ConfigUpdateAll, {}); // reload
                message.success("remove prefab " + targetKey + " succ");
            }

            this.setState({ selectedKey: "" })
            this.syncConfig()
        });
    }

    handleSearch = (selectedKeys, confirm, dataIndex) => {
        confirm();
        this.setState({ selectedKey: selectedKeys[0], SearchedColumn: dataIndex })
    };

    clickFmtBtn = (e) => {
        let old = this.state.code;
        this.setState({ code: formatText(old) })
    };

    showModal = () => {
        this.setState({ isModalVisible: true });
    };

    modalConfigChange = (e) => {
        this.setState({ newPrefabName: e.target.value });
    };

    modalHandleOk = () => {
        this.setState({ isModalVisible: false });

        let prefabName = this.state.newPrefabName
        if (prefabName === "") { return }
        prefabName = prefabName.toLowerCase()

        if (prefabName === "system" || prefabName === "global") {
            message.warning("System keywords that cannot be used!")
        }

        this.uploadPrefab(prefabName, "")
    };

    modalHandleCancel = () => {
        this.setState({ isModalVisible: false });
    };

    updateTags(name, tags) {
        /*
        var bots = this.state.Bots
        var tagSet = new Set()

        for (var i = 0; i < bots.length; i++) {
            if (bots[i].Name === name) {
                bots[i].Tags = tags   // update tags
            }

            if (bots[i].Tags) {
                for (var j = 0; j < bots[i].Tags.length; j++) {
                    tagSet.add(bots[i].Tags[j])
                }
            }
        }
        */
    }

    render() {

        const { isModalVisible } = this.state;

        const options = {
            mode: "text/x-lua",
            theme: localStorage.theme,
            lineNumbers: true,
        };

        // rowSelection object indicates the need for row selection
        const rowSelection = {
            onChange: (selectedRowKeys, selectedRows) => {
                console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
                this.setState({ code: selectedRows[0].code, selectedKey: selectedRows[0].key })
            },
            getCheckboxProps: (record) => ({
                // Column configuration not to be checked
                name: record.name,
            }),
        };

        return (

            <div>
                <ReflexContainer orientation="vertical">

                    <ReflexElement className="left-pane" flex={0.3} minSize="200">

                        <Table
                            columns={[
                                {
                                    title: 'Name',
                                    dataIndex: 'name',
                                    key: 'name',
                                    ...this.getColumnSearchProps('name')
                                },
                                {
                                    title: "Tags",
                                    dataIndex: "tags",
                                    key: "tags",
                                    render: (text, record) => (
                                        <HomeTagGroup record={record} onChange={(tags) => {
                                            this.updateTags(record.name, tags)
                                        }} ></HomeTagGroup>
                                    ),
                                },
                            ]}
                            onRow={(record) => {
                                return {
                                    onClick: (e) => {
                                        e.currentTarget.getElementsByClassName("ant-radio-wrapper")[0].click()
                                    },       // 点击行
                                };
                            }}
                            rowSelection={{
                                type: "radio",
                                ...rowSelection,
                            }}
                            dataSource={this.state.data}
                            pagination={{
                                pageSize: 100,
                            }}
                            scroll={{
                                y: document.body.clientHeight - 250,
                            }}

                        />

                        <Space size={20}>
                            <Button
                                onClick={this.onHandleAdd}
                                type="primary"
                                style={{
                                    marginBottom: 16,
                                }}
                            >
                                Add a prefab
                            </Button>

                            <Button
                                onClick={this.onHandleUpdate}
                                style={{
                                    marginBottom: 16,
                                }}
                            >
                                Update a prefab
                            </Button>

                            <Tooltip
                                placement="bottomLeft"
                                title={lanMap["app.prefab.remove.desc"][moment.locale()]}
                            >
                                <Popconfirm
                                    title={lanMap["app.prefab.remove.confirm"][moment.locale()]}
                                    onConfirm={this.onHandleRemove}
                                    onCancel={(e) => { }}
                                    okText="Yes"
                                    cancelText="No"
                                >
                                    <Button
                                        type="dashed"
                                        danger
                                        style={{
                                            marginBottom: 16,
                                        }}
                                    >
                                        Remove
                                    </Button>
                                </Popconfirm>
                            </Tooltip>
                        </Space>

                    </ReflexElement>

                    <ReflexSplitter propagate={true} />

                    <ReflexElement className="right-pane" flex={0.7} minSize="100">
                        <CodeMirror
                            value={this.state.code}
                            options={options}
                            onBeforeChange={this.onBeforeChange}
                        />
                        <Button onClick={this.clickFmtBtn}>Format script</Button>
                    </ReflexElement>

                </ReflexContainer>

                <Modal
                    visible={isModalVisible}
                    onOk={this.modalHandleOk}
                    onCancel={this.modalHandleCancel}
                >
                    <Input
                        placeholder={lanMap["app.prefab.add.placeholder"][moment.locale()]}
                        onChange={this.modalConfigChange}
                    />
                </Modal>
            </div>
        );
    }

}