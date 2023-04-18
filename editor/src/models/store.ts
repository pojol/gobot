import { createSlice, PayloadAction, configureStore } from '@reduxjs/toolkit';
import { combineReducers } from 'redux';


import prefabSlice from "@/models/prefab"
import treeSlice from "@/models/tree"
import debugInfoSlice from "@/models/debuginfo"
import configSlice from './config';
import resizeSlice from './resize';

const rootReducer = combineReducers({
    prefabSlice: prefabSlice.reducer,
    treeSlice: treeSlice.reducer,
    debugInfoSlice: debugInfoSlice.reducer,
    configSlice:configSlice.reducer,
    resizeSlice : resizeSlice.reducer,
  });
  

const store = configureStore ({
    reducer: rootReducer
});

export type RootState = ReturnType<typeof rootReducer>;
export default store;