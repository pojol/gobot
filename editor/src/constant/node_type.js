const NodeTy = {
  // 行为树根节点
  Root : "RootNode",

  // 脚本动作节点
  Action: "ActionNode",
  // 脚本条件节点
  Condition: "ConditionNode",
  // 脚本断言节点
  Assert: "AssertNode",

  // 循环次数节点
  Loop: "LoopNode",
  // 阻塞等待节点
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
