import React, { useState, useEffect } from 'react';

import { NodeTy } from "../../../constant/node_type";

import LoopTab from "./loop";
import WaitTab from "./wait";
import SequenceTab from "./sequence";

import { useSelector } from "react-redux";
import { RootState } from "@/models/store";
import ScriptTab from './script';

/// <reference path="node.d.ts" />

function GetPane(clickinfo: NodeClickInfo) {

  switch (clickinfo.type) {
    case NodeTy.Sequence:
    case NodeTy.Selector:
    case NodeTy.Parallel:
    case NodeTy.Root:
      return <SequenceTab />;
    case NodeTy.Wait:
      return <WaitTab />;
    case NodeTy.Loop:
      return <LoopTab />;
    default:
      return <ScriptTab />;
  }
}

export default function Nodes() {

  const { currentClickNode } = useSelector((state: RootState) => state.treeSlice);

  useEffect(() => {
    
  },[currentClickNode])

  return (
    <div>
      <GetPane {...currentClickNode} />
    </div>
  );
}
