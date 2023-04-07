
import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";


interface DebugInfoState {
    threadInfo : Array<ThreadInfo>
    metaInfo : string
}

const initialState : DebugInfoState = {
    metaInfo:"{}",
    threadInfo: new Array<ThreadInfo>()
}

const debugInfoSlice = createSlice({
    name: "debuginfo",
    initialState,
    reducers: {
        setDebugInfo(state, action: PayloadAction<DebugInfoState>) {
            state.metaInfo = action.payload.metaInfo
            state.threadInfo = action.payload.threadInfo
        }
    }
})

export const { setDebugInfo } = debugInfoSlice.actions;
export default debugInfoSlice;
