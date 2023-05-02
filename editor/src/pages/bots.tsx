import {
  Table,
  Space,
  InputNumber,
  Button,
  Upload,
  message,
  Popconfirm,
  Tooltip,
  Col,
  Row,
  Select,
} from "antd";
import React, { useState, useEffect } from 'react';
import {
  InboxOutlined,
  CloudDownloadOutlined,
  VerticalAlignBottomOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  CheckCircleTwoTone,
  CloseCircleTwoTone,
  ExclamationCircleTwoTone
} from "@ant-design/icons";

import { useLocation, history } from 'umi';
import { RootState } from '@/models/store';
import { connect, ConnectedProps } from 'react-redux';


import PubSub from "pubsub-js";
import Topic from "../constant/topic";
import Api from "../constant/api";
import { HomeTag } from "./tags/tags";
import { initTree } from "@/models/tree";

const { Post } = require("../utils/request");
const { LoadBehaviorWithBlob, LoadBehaviorWithFile } = require('../utils/parse');


const { Dragger } = Upload;
const { Option } = Select;

interface Record {
  name: string,
}

interface BotInfo {
  Tags: Array<string>,
  Name: string,
  Update: number,
  Status: string,
  Desc: string,
}

interface BotsProps extends PropsFromRedux { }

const Bots = (props: BotsProps) => {

  const [runs, setRuns] = useState(new Map<string, number>());
  const [batchLst, setBatchLst] = useState([]);

  // 表引用的数据
  const [botLst, setBotLst] = useState<{ name: string; key: string; update: string; status: string[]; tags: string[]; desc: string; }[]>([]);
  // 服务器存储的数据
  const [Bots, setBots] = useState<BotInfo[]>([]);
  // 可选的 tags
  const [selectedTags, setSelectedTags] = useState([]);
  // 当前选中的 tags
  const [currentSelectedTags, setCurrentSelectedTags] = useState<Array<string>>([]);
  // 选中的行
  const [selectedRows, setSelectedRows] = useState<Record[]>([]);

  const columns = [
    {
      title: "Behavior tree files",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Tags",
      dataIndex: "tags",
      key: "tags",
      render: (text: string, record: any) => (
        <HomeTag record={record} onChange={(tags: any) => {
          updateTags(record.name, tags)
        }} ></HomeTag>
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
      render: (text: string, record: any) => (
        <InputNumber
          min={0}
          max={100000}
          defaultValue={0}
          onChange={(e: any) => {
            var old = runs
            old.set(record.name, e)
            setRuns(old)
          }}
        ></InputNumber>
      ),
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: any, record: any) => (
        <>
          {status.map((s: string) => {
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
  ]

  useEffect(() => {
    refreshBotList();

    const token = PubSub.subscribe(Topic.BotsUpdate, () => {
      refreshBotList();
    });

    return () => {
      PubSub.unsubscribe(token);
    };
  }, []);

  useEffect(() => {
    fillBotList(Bots)
  }, [currentSelectedTags]);


  function fillBotList(bots: Array<BotInfo>) {
    var selectedTag = currentSelectedTags

    var intags = (tags: Array<string>) => {
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
        var tags = new Array<string>

        if (bots[i].Tags) {
          tags = bots[i].Tags
        }

        if (selectedTag.length > 0) {
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

      setBotLst(botlist);
    }

  }

  function updateTags(name: string, tags: Array<string>) {

    var bots = Bots
    var tagSet = new Set<string>()

    for (var i = 0; i < bots.length; i++) {
      if (bots[i].Name === name) {
        bots[i].Tags = tags   // update tags
        // 同步给服务器
        Post(localStorage.remoteAddr, Api.FileSetTags, {
          Name: name,
          NewTags: tags,
        }).then((json: any) => {
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
    setSelectedTags(children)
    setBots(bots)
    fillBotList(bots)
  }

  function updateAllTags(bots: Array<BotInfo>) {
    var tagSet = new Set<string>()

    for (var i = 0; i < bots.length; i++) {
      if (bots[i].Tags) {
        for (var j = 0; j < bots[i].Tags.length; j++) {
          tagSet.add(bots[i].Tags[j])
        }
      }
    }

    var children = []
    for (let tag of tagSet.keys()) {
      if (tag as string) {
        children.push(<Option key={tag} value={tag}>{tag}</Option>)
      }
    }

    setSelectedTags(children)
  }

  function refreshBotList() {
    setBatchLst([])

    console.info(localStorage.remoteAddr + "/" + Api.FileList)
    Post(localStorage.remoteAddr, Api.FileList, {}).then((json: any) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        if (json.Body.Bots) {
          setBots(json.Body.Bots)
          updateAllTags(json.Body.Bots)
          fillBotList(json.Body.Bots);
        }
      }
    });
  }

  const handleSelectChange = (tags: Array<string>) => {
    setCurrentSelectedTags(tags)
  }

  const handleBotLoad = (e: any) => {
    selectedRows
    if (selectedRows.length > 0) {
      var row = selectedRows[0]
      LoadBehaviorWithBlob(
        localStorage.remoteAddr,
        Api.FileGet,
        row.name
      ).then((file: any) => {
        LoadBehaviorWithFile(row.name, file.blob, (tree: any) => {
          props.dispatch(initTree(tree))
          history.push('/editor')
        });

      });
    }
  }

  const handleBotRun = (e: any) => {

    for (var i = 0; i < selectedRows.length; i++) {
      var row = selectedRows[i]

      var num = runs.get(row.name)
      if (num === undefined || num === 0) {
        Post(localStorage.remoteAddr, Api.BotRun, { Name: row.name }).then((json: any) => {
          if (json.Code !== 200) {
            message.error("running fail:" + String(json.Code) + " msg: " + json.Msg);
          } else {
            message.success("running succ");
          }
        })
        continue
      }

      Post(localStorage.remoteAddr, Api.BotCreateBatch, { Name: row.name, Num: num }).then((json: any) => {
        if (json.Code !== 200) {
          message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
        } else {
          message.success("batch run succ");
        }
      });
    }

  }

  const handleBotDelete = (e: any) => {
    for (var i = 0; i < selectedRows.length; i++) {
      var row = selectedRows[i]
      Post(localStorage.remoteAddr, Api.FileRemove, {
        Name: row.name,
      }).then((json: any) => {
        if (json.Code !== 200) {
          message.error(
            "run fail:" + String(json.Code) + " msg: " + json.Msg
          );
        } else {
          refreshBotList();
          message.success("bot delete succ");
        }
      });
    }
  }

  const handleBotDownload = (e: any) => {

    for (var i = 0; i < selectedRows.length; i++) {
      var row = selectedRows[i]

      LoadBehaviorWithBlob(
        localStorage.remoteAddr,
        Api.FileGet,
        row.name
      ).then((file: any) => {
        console.info("file =>", file)
        // 创建一个blob的对象，把Json转化为字符串作为我们的值
        if (window.navigator.msSaveOrOpenBlob) {
          navigator.msSaveBlob(file.blob, file.name)
        } else {

          // 上面这个是创建一个blob的对象连链接，
          // 创建一个链接元素，是属于 a 标签的链接元素，所以括号里才是a，
          var link = document.createElement("a");
          let body = document.querySelector("body")

          link.href = window.URL.createObjectURL(file.blob);
          link.download = file.name;

          // firefox
          link.style.display = "node"
          body.appendChild(link)

          // 使用js点击这个链接
          link.click();
          body.removeChild(link)

          window.URL.revokeObjectURL(link.href)
        }
      });
    }

  }

  var filepProps = {
    name: "file",
    multiple: true,
    action: localStorage.remoteAddr + "/" + Api.FileTxtUpload,
    onDrop(e: any) {
      console.log("Dropped files", e.dataTransfer.files);
    },
    onChange(e: any) {
      if (e.file.status === "done") {
        PubSub.publish(Topic.BotsUpdate, {})
      }
    }
  };

  return (
    <div >
      <Dragger {...filepProps}>
        <p className="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p className="ant-upload-text">
          {"drop"}
        </p>
      </Dragger>

      <div >
        <Row>
          <Col span={6}>
            <Select
              mode="multiple"
              allowClear
              style={{ width: '100%' }}
              placeholder={"select"}
              onChange={handleSelectChange}
            >
              {selectedTags}
            </Select>
          </Col>
          <Col span={6} offset={6}>
            <Space >
              <Tooltip
                placement="bottomLeft"
                title={"run"}
              >
                <Button icon={<PlayCircleOutlined />} onClick={handleBotRun}>
                  {"run"}
                </Button>
              </Tooltip>
              <Tooltip
                placement="bottomLeft"
                title={"load"}
              >
                <Button
                  icon={<CloudDownloadOutlined />}
                  onClick={handleBotLoad}
                >
                  {"load"}
                </Button>
              </Tooltip>
              <Tooltip
                placement="bottomLeft"
                title={"delete"}
              >
                <Popconfirm
                  title={"delete"}
                  onConfirm={handleBotDelete}
                  onCancel={(e) => { }}
                  okText="Yes"
                  cancelText="No"
                >
                  <Button icon={<DeleteOutlined />}> {"delete"}</Button>
                </Popconfirm>
              </Tooltip>
              <Tooltip
                placement="bottomLeft"
                title={"download"}
              >
                <Button
                  icon={<VerticalAlignBottomOutlined />}
                  onClick={handleBotDownload}
                >
                  {"download"}
                </Button>
              </Tooltip>
            </Space>
          </Col>
        </Row>

      </div>
      <Table
        rowSelection={{
          type: "checkbox",
          ...{
            onChange: (selectedRowKeys: any, selectedRows: any) => {
              setSelectedRows(selectedRows)
              console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
            }
          },
        }}
        columns={columns}
        dataSource={botLst}
        onRow={(record) => {
          return {
            onClick: (e) => {
              e.currentTarget.getElementsByClassName("ant-checkbox-wrapper")[0].click()
            },       // 点击行
          };
        }}
        pagination={{
          position: ['bottomCenter'],
        }}
      />

    </div>
  );
}

const mapStateToProps = (state: RootState) => ({
});

const connector = connect(mapStateToProps);
type PropsFromRedux = ConnectedProps<typeof connector>;

export default connector(Bots);