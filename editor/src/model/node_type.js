const NodeTy = {
  Root : "RootNode",

  Action: "ActionNode",
  Condition: "ConditionNode",
  Assert: "AssertNode",
  Loop: "LoopNode",
  Wait: "WaitNode",

  Selector: "SelectorNode",
  Sequence: "SequenceNode",
};

function IsScriptNode(ty) {
  if (ty === NodeTy.Condition || ty === NodeTy.Action || ty === NodeTy.Assert) {
    return true;
  } else {
    return false;
  }
}

export { NodeTy, IsScriptNode };
