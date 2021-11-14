import {
  Table,
  Tag,
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
import * as React from "react";
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
import { NodeTy, IsScriptNode } from "../model/node_type";
import "./home.css";
import { SaveAs } from "../utils/file";


const { Dragger } = Upload;
const { Option } = Select;

function GetBehaviorBlob(url, methon, name) {
  return new Promise(function (resolve, reject) {
    fetch(url + methon, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: JSON.stringify({ Name: name }),
    })
      .then((response) => {
        if (response.ok) {
          return response.blob();
        } else {
          reject({ status: response.status });
        }
      })
      .then((response) => {
        resolve(response);
      })
      .catch((err) => {
        reject({ status: -1 });
      });
  });
}

function getValueByElement(elem, tag) {
  for (var i = 0; i < elem.childNodes.length; i++) {
    if (elem.childNodes[i].nodeName === tag) {
      if (elem.childNodes[i].childNodes.length === 0) {
        return ""
      } else {
        return elem.childNodes[i].childNodes[0].nodeValue;
      }
    }
  }
  return undefined;
}

function parseChildren(xmlnode, children) {
  var nod = {};

  nod.id = xmlnode.getElementsByTagName("id")[0].childNodes[0].nodeValue;
  nod.ty = xmlnode.getElementsByTagName("ty")[0].childNodes[0].nodeValue;

  if (nod.ty === NodeTy.Loop) {
    nod.loop = getValueByElement(xmlnode, "loop")
  } else if (nod.ty === NodeTy.Wait) {
    nod.wait = getValueByElement(xmlnode, "wait")
  } else if (IsScriptNode(nod.ty)) {
    nod.code = getValueByElement(xmlnode, "code");
    nod.alias = getValueByElement(xmlnode, "alias");
  }

  nod.pos = {
    x: parseInt(
      xmlnode.getElementsByTagName("pos")[0].getElementsByTagName("x")[0]
        .childNodes[0].nodeValue
    ),
    y: parseInt(
      xmlnode.getElementsByTagName("pos")[0].getElementsByTagName("y")[0]
        .childNodes[0].nodeValue
    ),
  };

  nod.children = [];
  children.push(nod);

  for (var i = 0; i < xmlnode.childNodes.length; i++) {
    if (xmlnode.childNodes[i].nodeName === "children") {
      parseChildren(xmlnode.childNodes[i], nod.children);
    }
  }
}

function LoadFile(name, blob) {
  let reader = new FileReader();
  reader.onload = function (ev) {
    var context = reader.result;
    try {
      let parser = new DOMParser();
      let xmlDoc = parser.parseFromString(context, "text/xml");

      let tree = {};
      var root = xmlDoc.getElementsByTagName("behavior")[0];
      if (root) {
        tree.id = root.getElementsByTagName("id")[0].childNodes[0].nodeValue;
        tree.ty = root.getElementsByTagName("ty")[0].childNodes[0].nodeValue;
        tree.pos = {
          x: parseInt(
            root.getElementsByTagName("pos")[0].getElementsByTagName("x")[0]
              .childNodes[0].nodeValue
          ),
          y: parseInt(
            root.getElementsByTagName("pos")[0].getElementsByTagName("y")[0]
              .childNodes[0].nodeValue
          ),
        };
        tree.children = [];
        if (root.getElementsByTagName("children")[0].hasChildNodes()) {
          parseChildren(
            root.getElementsByTagName("children")[0],
            tree.children
          );
        }
      }

      PubSub.publish(Topic.FileLoad, {
        Name: name,
        Tree: tree,
      });
    } catch (err) {
      console.info(err)
      message.warning("文件解析失败");
    }
  };

  reader.readAsText(blob);
}

export default class BotList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      searchText: "",
      searchedColumn: "",
      runs: {},
      selectedTags: [],
      selectedRows: [],
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
          render: tags => (
            <>
              {tags.map(tag => {
                return (
                  <Tag key={tag}>
                    {tag}
                  </Tag>
                );
              })}
            </>
          ),

        },
        {
          title: "UpdateTime",
          dataIndex: "update",
          key: "update",
          sorter: (a, b) => a.update > b.update,
        },
        {
          title: "Number of runs",
          dataIndex: "num",
          key: "num",
          render: (text, record) => (
            <InputNumber
              min={0}
              max={1000}
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
      botLst: [],
      batchLst: [],
    };
  }

  componentDidMount() {
    PubSub.subscribe(Topic.BotsUpdate, (topic, info) => {
      this.refreshBotList();
    });

    this.refreshBotList();
  }

  fillBotList(lst) {
    if (lst) {
      var botlist = [];
      var tags = ["aa", "bb"]

      for (var i = 0; i < lst.length; i++) {
        var _upt = new Date(lst[i].Update * 1000);
        var _upts = _upt.toLocaleDateString() + " " + _upt.toLocaleTimeString();
        botlist.push({
          name: lst[i].Name,
          key: lst[i].Name,
          update: _upts,
          status: [lst[i].Status],
          tags: ["aa", "bb"],
          desc: lst[i].Desc
        });
      }

      var children = []
      for (i = 0; i < tags.length; i++) {
        children.push(<Option key={tags[i]} value={tags[i]}>{tags[i]}</Option>)
      }

      this.setState({ botLst: botlist });
      this.setState({ selectedTags: children })
    }
  }

  refreshBotList() {
    this.setState({ botLst: [] });
    this.setState({ batchLst: [] });

    Post(window.remote, Api.FileList, {}).then((json) => {
      if (json.Code !== 200) {
        message.error("run fail:" + String(json.Code) + " msg: " + json.Msg);
      } else {
        this.fillBotList(json.Body.Bots);
      }
    });
  }

  onLoadClick = (key) => {
    console.info(key);
  };

  getColumnSearchProps = (dataIndex) => ({
    filterDropdown: ({
      setSelectedKeys,
      selectedKeys,
      confirm,
      clearFilters,
    }) => (
      <div style={{ padding: 8 }}>
        <Input
          ref={(node) => {
            this.searchInput = node;
          }}
          placeholder={`Search ${dataIndex}`}
          value={selectedKeys[0]}
          onChange={(e) =>
            setSelectedKeys(e.target.value ? [e.target.value] : [])
          }
          onPressEnter={() =>
            this.handleSearch(selectedKeys, confirm, dataIndex)
          }
          style={{ marginBottom: 8, display: "block" }}
        />
        <Space>
          <Button
            type="primary"
            onClick={() => this.handleSearch(selectedKeys, confirm, dataIndex)}
            icon={<SearchOutlined />}
            size="small"
            style={{ width: 90 }}
          >
            Search
          </Button>
          <Button
            onClick={() => this.handleReset(clearFilters)}
            size="small"
            style={{ width: 90 }}
          >
            Reset
          </Button>
        </Space>
      </div>
    ),
    filterIcon: (filtered) => (
      <SearchOutlined style={{ color: filtered ? "#1890ff" : undefined }} />
    ),
    onFilter: (value, record) =>
      record[dataIndex]
        ? record[dataIndex]
          .toString()
          .toLowerCase()
          .includes(value.toLowerCase())
        : "",
    onFilterDropdownVisibleChange: (visible) => {
      if (visible) {
        setTimeout(() => this.searchInput.select(), 100);
      }
    },
    render: (text) =>
      this.state.searchedColumn === dataIndex ? (
        <Highlighter
          highlightStyle={{ backgroundColor: "#ffc069", padding: 0 }}
          searchWords={[this.state.searchText]}
          autoEscape
          textToHighlight={text ? text.toString() : ""}
        />
      ) : (
        text
      ),
  });

  handleSearch = (selectedKeys, confirm, dataIndex) => {
    confirm();
    this.setState({
      searchText: selectedKeys[0],
      searchedColumn: dataIndex,
    });
  };

  handleReset = (clearFilters) => {
    clearFilters();
    this.setState({ searchText: "" });
  };

  uploadOnChange = (info) => {
    const { status } = info.file;
    if (status !== "uploading") {
      console.log(info.file, info.fileList);
    }
    if (status === "done") {
      message.success(`${info.file.name} file uploaded successfully.`);
      this.refreshBotList();
    } else if (status === "error") {
      message.error(`${info.file.name} file upload failed.`);
    }
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
      GetBehaviorBlob(
        window.remote,
        Api.FileGet,
        row.name
      ).then((blob) => {
        LoadFile(row.name, blob);
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

      GetBehaviorBlob(
        window.remote,
        Api.FileGet,
        row.name
      ).then((blob) => {
        // 创建一个blob的对象，把Json转化为字符串作为我们的值
        SaveAs(blob, row.name)
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

        <Table rowSelection={{
          type: "checkbox",
          ...this.rowSelection,
        }}
          columns={this.state.columns} dataSource={this.state.botLst} />

      </div>
    );
  }
}
