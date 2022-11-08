import * as React from "react";
import { Space, Table, Tag, Button } from 'antd';

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/solarized.css";
import "codemirror/theme/abcdef.css";
import "codemirror/theme/ayu-dark.css";
import "codemirror/theme/yonce.css";
import "codemirror/theme/neo.css";
import "codemirror/theme/zenburn.css";
import "codemirror/mode/lua/lua";
import { SwatchesPicker } from "react-color";

// You will need to import the styles separately
// You probably want to do this just once during the bootstrapping phase of your application.
import "react-reflex/styles.css";

// then you can import the components
import { ReflexContainer, ReflexSplitter, ReflexElement } from "react-reflex";
import { append } from "@antv/x6/lib/util/dom/elem";


const columns = [
    {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
        render: (text) => <a>{text}</a>,
    }
];


export default class BotPrefab extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: [
                {
                    key: '1',
                    name: 'John Brown',
                },
                {
                    key: '2',
                    name: 'Jim Green',
                },
                {
                    key: '3',
                    name: 'Joe Black',
                }
            ]
        };
    }

    componentDidMount() {

        let dat = []
        for (var i = 0; i < 1000; i++) {
            dat.push({ key: i.toString(), name: i.toString() })
        }

        console.info("data", dat)
        this.setState({ data: dat })
    }

    onBeforeChange = (editor, data, value) => {

    };

    onHandleAdd = () => {

    }

    render() {

        const { } = this.state;

        const options = {
            mode: "text/x-lua",
            theme: localStorage.theme,
            lineNumbers: true,
        };

        return (
            <div>

                <ReflexContainer orientation="vertical">

                    <ReflexElement className="left-pane" flex={0.3} minSize="200">
                        <Button
                            onClick={this.onHandleAdd}
                            type="primary"
                            style={{
                                marginBottom: 16,
                            }}
                        >
                            Add a prefab
                        </Button>
                        <Table
                            columns={columns}
                            dataSource={this.state.data}
                            pagination={{
                                pageSize: 100,
                            }}
                            scroll={{
                                y: document.body.clientHeight - 250,
                            }}

                        />
                    </ReflexElement>

                    <ReflexSplitter propagate={true} />

                    <ReflexElement className="right-pane" flex={0.7} minSize="100">
                        <CodeMirror
                            value={""}
                            options={options}
                            onBeforeChange={this.onBeforeChange}
                        />
                    </ReflexElement>

                </ReflexContainer>


            </div>
        );
    }

}