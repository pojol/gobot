import { createSlice, PayloadAction, configureStore } from '@reduxjs/toolkit';
import { combineReducers } from 'redux';


import prefabSlice from "@/models/mstore/prefab"
import treeSlice from "@/models/mstore/tree"
import debugInfoSlice from "@/models/mstore/debuginfo"

const rootReducer = combineReducers({
    prefabSlice: prefabSlice.reducer,
    treeSlice: treeSlice.reducer,
    debugInfoSlice: debugInfoSlice.reducer,
  });
  

const store = configureStore ({
    reducer: rootReducer
});

export type RootState = ReturnType<typeof rootReducer>;
export default store;