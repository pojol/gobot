import React, { useState, useEffect } from "react";
import {
    InputNumber,
    Row,
    Col,
    Button,
    message,
    Slider,
    Space,
    Input,
} from "antd";
import {  useSelector } from 'react-redux';
import { RootState } from "@/models/store";

/// <reference path="node.d.ts" />

export default function SequenceTab() {
    const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);

    return (
        <div>
            <Button type="dashed">{currentClickNode.id}</Button>
        </div>
    );
}
