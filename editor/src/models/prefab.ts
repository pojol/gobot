import { enableMapSet } from 'immer';
import { createSlice, PayloadAction } from '@reduxjs/toolkit';
const { Post, PostGetBlob } = require("../utils/request");
import Api from '@/constant/api';
import PubSub from "pubsub-js";
import Topic from "../constant/topic";
import store from './store';

enableMapSet(); // 启用 Map 对象的支持

const initmap = (): Array<PrefabInfo> => {
    let map = new Array()
    if (localStorage.remoteAddr === undefined || localStorage.remoteAddr === "") {
        // 直接返回空的map
        return map
    }

    Post(localStorage.remoteAddr, Api.PrefabList, {}).then((json: any) => {
        if (json.Code == 200) {
            let lst = json.Body.Lst;
            var counter = 0;

            if (lst) {
                lst.forEach(function (element: any) {
                    PostGetBlob(localStorage.remoteAddr, Api.PrefabGet, element.name).then(
                        (file: any) => {
                            let reader = new FileReader();
                            reader.onload = function (ev) {
                                let pi: PrefabInfo = {
                                    name: element.name,
                                    tags: element.tags,
                                    code: String(reader.result),
                                }
    
                                store.dispatch(addItem({ key: element.name, value: pi }))
    
                                counter++;
                                if (counter === lst.length) {
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

    return map
}

interface PrefabState {
    pmap: Array<PrefabInfo>
}


const initialState: PrefabState = {
    pmap: initmap(),
};


const prefabSlice = createSlice({
    name: 'prefab',
    initialState,
    reducers: {
        addItem(state, action: PayloadAction<{ key: string; value: PrefabInfo }>) {
            const { key, value } = action.payload;
            state.pmap.push(value)
        },
        removeItem(state, action: PayloadAction<string>) {
            const key = action.payload;
            state.pmap = state.pmap.filter(item => item.name !== key)
        },
        cleanItems(state, action: PayloadAction<void>) {
            state.pmap = new Array()
        }
    },
});

export const { addItem, removeItem, cleanItems } = prefabSlice.actions;
export default prefabSlice;