import { createAction, createSlice, PayloadAction } from "@reduxjs/toolkit";

interface ResizeState {
    editFlex: number
    graphFlex: number
}

const initialState: ResizeState = {
    graphFlex: 0.6,  // graph 在屏幕中的占比
    editFlex: 0.4   // edit（编辑窗）在屏幕中的占比
}

const resizeSlice = createSlice({
    name: "debuginfo",
    initialState,
    reducers: {
        setGraphFlex(state, action: PayloadAction<number>) {
            state.graphFlex = action.payload
        },
        setEditFlex(state, action: PayloadAction<number>) {
            state.editFlex = action.payload
        },
    }
})

export const { setGraphFlex, setEditFlex } = resizeSlice.actions;
export default resizeSlice;