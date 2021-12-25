import {
  Table,
  Space,
  InputNumber,
  Button,
  Upload,
  message,
  Input,
  Popconfirm,
  Tooltip,
  Col,
  Row,
  Select,
} from "antd";
import React, { } from 'react';
import {
  InboxOutlined,
  CloudDownloadOutlined,
  SearchOutlined,
  VerticalAlignBottomOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  CheckCircleTwoTone,
  CloseCircleTwoTone,
  ExclamationCircleTwoTone
} from "@ant-design/icons";
import Highlighter from "react-highlight-words";
import PubSub from "pubsub-js";
import Topic from "../model/topic";
import { Post } from "../model/request";
import Api from "../model/api";
import "./home.css";
import { SaveAs } from "../utils/file";
import { LoadBehaviorWithBlob, LoadBehaviorWithFile } from "../utils/tree";
import HomeTagGroup from "./home_tags";
import { set } from "@antv/util";

const { Dragger } = Upload;
const { Option } = Select;

export default class BotList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      runs: {},
      columns: [
        {
          title: "Bot behavior file",
          dataIndex: "name",
          key: "name",
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
        {
          title: "UpdateTime",
          dataIndex: "update",
          key: "update",
          //sorter: (a, b) => a.update > b.update,
        },
        {
          title: "Number of runs",
          dataIndex: "num",
          key: "num",
          render: (text, record) => (
            <InputNumber
              min={0}
              max={100000}
              defaultValue={0}
              onChange={(e) => {
                var old = this.state.runs
                old[record.name] = e
                this.setState({ runs: old })
              }}
            ></InputNumber>
          ),
        },
        {
          title: "Status",
          dataIndex: "status",
          key: "status",
          render: (status, record) => (
            <>
              {status.map(s => {
                if (s === 'succ') {
                  return <CheckCircleTwoTone twoToneColor="#52c41a" />
                } else if (s === 'fail') {
                  return <CloseCircleTwoTone twoToneColor="#eb2f96" />
                } else {
                  return <ExclamationCircleTwoTone twoToneColor='#adb5bd' />
                }
              })}
            </>
          ),
        },
        {
          title: "Desc",
          dataIndex: "desc",
          key: "desc",
        }
      ],
      Bots: [],

      batchLst: [],
      botLst: [],               // 显示的 botlist

      selectedTags: [],         // 可选的 tags
      currentSelectedTags: [], // 当前选中的 tags

      selectedRows: [],         // 选中的行
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.BotsUpdate, (topic, info) => {
      this.refreshBotList();
    });

  }

  fillBotList() {

    var bots = this.state.Bots
    var selectedTag = this.state.currentSelectedTags

    var intags = (tags) => {
      for (var i = 0; i < selectedTag.length; i++) {
        for (var j = 0; j < tags.length; j++) {
          if (selectedTag[i] === tags[j]) {
            return true
          }
        }
      }
      return false
    }

    if (bots.length > 0) {
      var botlist = [];
      
      for (var i = 0; i < bots.length; i++) {
        var tags = []

        if (bots[i].Tags) {
          tags = bots[i].Tags
        }

        if (selectedTag.length > 0) {
          console.info("select filter tags", selectedTag, "bot tags", bots[i].Tags)
          if (tags.length > 0) {
            if (!intags(tags)) {
              continue
            }
          } else {
            continue
          }
        }

        var _upt = new Date(bots[i].Update * 1000);
        var _upts = _upt.toLocaleDateString() + " " + _upt.toLocaleTimeString();
        botlist.push({
          name: bots[i].Name,
          key: bots[i].Name,
          update: _upts,
          status: [bots[i].Status],
          tags: tags,
          desc: bots[i].Desc
        });
      }

      this.setState({ botLst: botlist });
    }

  }

  updateTags(name, tags) {
    var bots = this.state.Bots
    var tagSet = new Set()

    console.info("update tags", name, tags)

    for (var i = 0; i < bots.length; i++) {
      if (bots[i].Name === name) {
        bots[i].Tags = tags   // update tags
        // 同步给服务器
        Post(window.remote, Api.FileSetTags, {
          Name : name,
          NewTags : tags,
        }).then((json)=>{
          if (json.Code !== 200) {
            message.error("updaet tags fail:" + String(json.Code) + " msg: " + json.Msg);
          } else {
            message.success("update tags succ!")
          }
        })
      }

      if (bots[i].Tags) {
        for (var j = 0; j < bots[i].Tags.length; j++) {
          console.info("add tag", bots[i].Tags[j])
          tagSet.add(bots[i].Tags[j])
        }
      }

    }

    var children = []
    for (let tag of tagSet.keys()) {
      children.push(<Option key={tag} value={tag}>{tag}</Option>)
    }

    // refresh tags
    console.info("refresh tags", children)
    this.setState({ selectedTags: children })

    this.setState({ Bots: bots })
    this.fillBotList()
  }

  updateAllTags() {
    var bots = this.state.Bots
    var tagSet = new Set()

    for (var i = 0; i < bots.length; i++) {
      if (bots[i].Tags) {
        for (var j = 0; j < bots[i].Tags.length; j++) {
          tagSet.add(bots[i].Tags[j])
        }
      }
    }

    var children = []
    for (let tag of tagSet.keys()) {
      children.push(<Option key={tag} value={tag}>{tag}</Option>)
    }

    console.info("selected tags", children)
    this.setState({ selectedTags: children })
  }

  refreshBotList() {
    this.setState({ botLst: [] });
    this.setState({ batchLst: [] });

    Post(window.remote, Api.FileList, {}).then((json) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        console.info("refresh bots", json.Body.Bots)
        if (json.Body.Bots) {
          this.setState({ Bots: json.Body.Bots }, ()=>{
            this.updateAllTags()
            this.fillBotList();
          })
        }
      }
    });
  }

  onLoadClick = (key) => {
    console.info(key);
  };


  refreshBatchInfo(name, cnt) {
    var flag = false;
    var old = this.state.batchLst;

    for (var i = 0; i < old.length; i++) {
      if (old[i].name === name) {
        old[i].cnt = cnt;
        flag = true;
      }
    }

    if (!flag) {
      old.push({
        name: name,
        cnt: cnt,
      });
    }

    this.setState({ batchLst: old });
  }

  handleSelectChange = (tags) => {

    console.info("refresh bot lst", tags)

    this.setState({ currentSelectedTags: tags }, () => {
      this.fillBotList()
    })

  }


  rowSelection = {
    onChange: (selectedRowKeys, selectedRows) => {
      this.setState({ selectedRows: selectedRows })
      console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
    }
  }

  handleBotLoad = e => {

    if (this.state.selectedRows.length > 0) {
      var row = this.state.selectedRows[0]
      LoadBehaviorWithBlob(
        window.remote,
        Api.FileGet,
        row.name
      ).then((file) => {
        var tree = LoadBehaviorWithFile(row.name, file.blob);
        if (tree !== null) {
          PubSub.publish(Topic.FileLoad, {
            Name: row.name,
            Tree: tree,
          });
        } else {
          message.warning("文件解析失败");
        }
      });
    }


  }

  handleBotRun = e => {

    for (var i = 0; i < this.state.selectedRows.length; i++) {
      var row = this.state.selectedRows[i]

      var num = this.state.runs[row.name]
      if (num === undefined || num === 0) {
        message.warn("Please set the number of bot runs " + row.name)
        continue
      }

      Post(window.remote, Api.BotCreate, { Name: row.name, Num: num }).then((json) => {
        if (json.Code !== 200) {
          message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
        } else {
          message.success("batch run succ");
        }
      });
    }


  }

  handleBotDelete = e => {
    for (var i = 0; i < this.state.selectedRows.length; i++) {
      var row = this.state.selectedRows[i]
      Post(window.remote, Api.FileRemove, {
        Name: row.name,
      }).then((json) => {
        if (json.Code !== 200) {
          message.error(
            "run fail:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          this.refreshBotList();
          message.success("bot delete succ");
        }
      });
    }
  }

  handleBotDownload = e => {

    for (var i = 0; i < this.state.selectedRows.length; i++) {
      var row = this.state.selectedRows[i]

      LoadBehaviorWithBlob(
        window.remote,
        Api.FileGet,
        row.name
      ).then((file) => {
        // 创建一个blob的对象，把Json转化为字符串作为我们的值
        SaveAs(file.blob, file.name)
      });
    }

  }

  render() {
    var filepProps = {
      name: "file",
      multiple: true,
      action: window.remote + Api.FileTxtUpload,
      onChange: this.uploadOnChange,
      onDrop(e) {
        console.log("Dropped files", e.dataTransfer.files);
      },
    };

    return (
      <div >
        <Dragger {...filepProps}>
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">
            Click or drag file (*.xml) to this area to upload
          </p>
        </Dragger>

        <div >
          <Row>
            <Col span={6}>
              <Select
                mode="multiple"
                allowClear
                style={{ width: '100%' }}
                placeholder="Filter by tags"
                onChange={this.handleSelectChange}
              >
                {this.state.selectedTags}
              </Select>
            </Col>
            <Col span={6} offset={6}>
              <Space >
                <Tooltip
                  placement="topLeft"
                  title="Drive a specified number of robots"
                >
                  <Button icon={<PlayCircleOutlined />} onClick={this.handleBotRun}>
                    Run
                  </Button>
                </Tooltip>
                <Tooltip
                  placement="topLeft"
                  title="Load the behavior file to the local for editing"
                >
                  <Button
                    icon={<CloudDownloadOutlined />}
                    onClick={this.handleBotLoad}
                  >
                    Load
                  </Button>
                </Tooltip>
                <Tooltip
                  placement="topLeft"
                  title="Delete the behavior file from the database"
                >
                  <Popconfirm
                    title="Are you sure to delete this bot?"
                    onConfirm={this.handleBotDelete}
                    onCancel={(e) => { }}
                    okText="Yes"
                    cancelText="No"
                  >
                    <Button icon={<DeleteOutlined />}>Delete</Button>
                  </Popconfirm>
                </Tooltip>
                <Tooltip
                  placement="topLeft"
                  title="Save the current behavior tree file to the local"
                >
                  <Button
                    icon={<VerticalAlignBottomOutlined />}
                    onClick={this.handleBotDownload}
                  >
                    Download
                  </Button>
                </Tooltip>
              </Space>
            </Col>
          </Row>


        </div>

        <Table
          rowSelection={{
            type: "checkbox",
            ...this.rowSelection,
          }}
          columns={this.state.columns}
          dataSource={this.state.botLst} />

      </div>
    );
  }
}
