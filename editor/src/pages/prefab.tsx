import React, { useState, useEffect } from 'react';
import {
  Space, Table, Button, message, Input, Popconfirm,
  Tooltip, Modal
} from 'antd';
import PubSub from "pubsub-js";
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


import Topic from "../constant/topic";
import Api from "../constant/api";
import { HomeTag } from "./tags/tags";
import type { ColumnsType, ColumnType } from 'antd/es/table';

import { connect, ConnectedProps } from 'react-redux';
import { addItem, removeItem, cleanItems } from '../models/prefab';
import { RootState } from '@/models/store';

const { Post, PostGetBlob, PostBlob } = require("../utils/request");

// You will need to import the styles separately
// You probably want to do this just once during the bootstrapping phase of your application.
import "react-reflex/styles.css";

// then you can import the components
import { ReflexContainer, ReflexSplitter, ReflexElement } from "react-reflex";

interface DataType {
  key: string;
  name: string;
  tags: Array<string>,
  code: string,
}

type DataIndex = keyof DataType;

interface PrefabProps extends PropsFromRedux { }

const Prefab = (props: PrefabProps) => {

  const [data, setData] = useState<DataType[]>([]);
  const [code, setCode] = useState<string>("");
  const [selectedKey, setSelectedKey] = useState<string>("");
  const [isModalVisible, setIsModalVisible] = useState<boolean>(false);
  const [newPrefabName, setNewPrefabName] = useState<string>("");
  const [searchedColumn, setSearchedColumn] = useState<DataIndex>();
  const [editor, setEditor] = useState<any>({});

  let obj: any = null

  useEffect(() => {
    syncConfig();
  }, []);


  const getColumnSearchProps = (dataIndex: DataIndex): ColumnType<DataType> => ({
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
          onPressEnter={() => handleSearch(selectedKeys, confirm, dataIndex)}
          style={{
            marginBottom: 8,
            display: 'block',
          }}
        />
      </div>
    ),
    filterIcon: (filtered: any) => (
      <FilterTwoTone />
    ),
    onFilter: (value, record) =>
      record[dataIndex].toString().toLowerCase().includes((value as string).toLowerCase()),
    render: (text: string) => <a>{text}</a>
  })

  const syncConfig = () => {

    props.dispatch(cleanItems())

    Post(localStorage.remoteAddr, Api.PrefabList, {}).then((json: any) => {
      if (json.Code !== 200) {
        message.error(
          "get config list fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        console.info("prefab lst", json.Body.Lst)
        let lst = json.Body.Lst;
        var counter = 0;
        let dat = new Array<DataType>;

        let callback = () => {
          dat.sort(function (a, b) {
            if (a.name < b.name) { return -1; }
            if (a.name > b.name) { return 1; }
            return 0;
          });

          setData(dat)
        };

        if (lst) {
          lst.forEach(function (element: any) {
            PostGetBlob(localStorage.remoteAddr, Api.PrefabGet, element.name).then(
              (file: any) => {
                let reader = new FileReader();
                reader.onload = function (ev) {

                  dat.push({ key: element.name, name: element.name, tags: element.tags, code: String(reader.result), })

                  let pi: PrefabInfo = {
                    name: element.name,
                    tags: element.tags,
                    code: String(reader.result),
                  }
                  props.dispatch(addItem({ key: element.name, value: pi }))

                  counter++;
                  if (counter === lst.length) {
                    callback()
                    PubSub.publish(Topic.PrefabUpdateAll, {})
                  }
                };

                reader.readAsText(file.blob);
              }
            );
          });
        }
      }
    });
  }

  const onBeforeChange = (editor: any, data: any, value: string) => {
    setCode(value)
  };

  const uploadPrefab = (name: string, code: string) => {
    var blob = new Blob([code], {
      type: "application/json",
    });

    PostBlob(
      localStorage.remoteAddr,
      Api.PrefabUpload,
      name,
      blob
    ).then((json: any) => {
      if (json.Code !== 200) {
        message.error(
          "upload fail:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        message.success(name + " upload succ");
      }

      syncConfig()
    });
  }

  const onHandleAdd = () => {
    setIsModalVisible(true)
  }

  const onHandleUpdate = () => {
    let targetKey = selectedKey
    if (targetKey === "") {
      return
    }

    uploadPrefab(targetKey, code)
  }

  const onHandleRemove = () => {
    let targetKey = selectedKey
    if (targetKey === "") {
      return
    }

    Post(localStorage.remoteAddr, Api.PrefabRemove, {
      Name: targetKey,
    }).then((json: any) => {
      if (json.Code !== 200) {
        message.error(
          "remove prefab err:" + String(json.Code) + " msg: " + json.Msg
        );
      } else {
        props.dispatch(removeItem(targetKey))
        PubSub.publish(Topic.PrefabUpdateAll, {}); // reload
        message.success("remove prefab " + targetKey + " succ");
      }

      setSelectedKey("")
      syncConfig()
    });
  }

  const handleSearch = (selectedKeys: any, confirm: any, dataIndex: any) => {
    confirm();
    setSelectedKey(selectedKeys[0])
    setSearchedColumn(dataIndex)
  };

  const showModal = () => {
    setIsModalVisible(true)
  };

  const modalConfigChange = (e: any) => {
    setNewPrefabName(e.target.value)
  };

  const modalHandleOk = () => {
    setIsModalVisible(false)

    let prefabName = newPrefabName
    if (prefabName === "") { return }
    prefabName = prefabName.toLowerCase()

    if (prefabName === "system" || prefabName === "global") {
      message.warning("System keywords that cannot be used!")
    }

    uploadPrefab(prefabName, "")
  };

  const modalHandleCancel = () => {
    setIsModalVisible(false)
  };

  const updateTags = (name: string, tags: Array<string>) => {
    const newData = [...data]; // create a new copy of the data array
    var tagSet = new Set()

    for (var i = 0; i < newData.length; i++) {
      if (newData[i].name === name) {
        console.info("update tags", name, tags)
        newData[i].tags = tags  // modify the copy of the data array
        Post(localStorage.remoteAddr, Api.PrefabSetTags, { name: name, tags: tags }).then((json: any) => {
          if (json.Code !== 200) {
            message.error("updaet tags fail:" + String(json.Code) + " msg: " + json.Msg);
          } else {
            message.success("update tags succ!")
          }
        })
      }

      if (newData[i].tags) {
        for (var j = 0; j < newData[i].tags.length; j++) {
          tagSet.add(newData[i].tags[j])
        }
      }
    }

    setData(newData)
  }


  const onDidMount = (editor: any) => {
    editor.setSize(undefined, document.documentElement.clientHeight - 120)
    setEditor(editor)
  };

  const options = {
    mode: "text/x-lua",
    theme: localStorage.codeboxTheme,
    lineNumbers: true,
    indentUnit: 4,
  };

  // rowSelection object indicates the need for row selection
  const rowSelection = {
    onChange: (selectedRowKeys: any, selectedRows: any) => {
      console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
      setCode(selectedRows[0].code)
    },
    getCheckboxProps: (record: any) => ({
      // Column configuration not to be checked
      name: record.name,
    }),
  };

  return (

    <div>
      <ReflexContainer orientation="vertical">

        <ReflexElement className="left-pane" flex={0.3} minSize={200}>

          <Table
            columns={[
              {
                title: 'Name',
                dataIndex: 'name',
                key: 'name',
                ...getColumnSearchProps('name')
              },
              {
                title: "Tags",
                dataIndex: "tags",
                key: "tags",
                render: (text: string, record: any) => (
                  <HomeTag record={record} onChange={(tags) => {
                    console.info("tags", tags)
                    updateTags(record.name, tags)
                  }} ></HomeTag>
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
            dataSource={data}
            pagination={{
              pageSize: 100,
            }}
            scroll={{
              y: document.documentElement.clientHeight - 270,
            }}

          />

          <Space size={20}>
            <Button
              onClick={onHandleAdd}
              type="primary"
              style={{
                marginBottom: 16,
              }}
            >
              Add a prefab
            </Button>

            <Button
              onClick={onHandleUpdate}
              style={{
                marginBottom: 16,
              }}
            >
              Update a prefab
            </Button>

            <Tooltip
              placement="bottomLeft"
              title={"app.prefab.remove.desc"}
            >
              <Popconfirm
                title={"app.prefab.remove.confirm"}
                onConfirm={onHandleRemove}
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

        <ReflexElement className="right-pane" flex={0.7} minSize={100}>
          <CodeMirror
            value={code}
            options={options}
            onBeforeChange={onBeforeChange}
            editorDidMount={onDidMount}
          />
        </ReflexElement>

      </ReflexContainer>

      <Modal
        open={isModalVisible}
        onOk={modalHandleOk}
        onCancel={modalHandleCancel}
      >
        <Input
          placeholder={"app.prefab.add.placeholder"}
          onChange={modalConfigChange}
        />
      </Modal>
    </div>
  );
}

const mapStateToProps = (state: RootState) => ({
  prefabMap: state.prefabSlice.pmap
});

const connector = connect(mapStateToProps);
type PropsFromRedux = ConnectedProps<typeof connector>;

export default connector(Prefab);