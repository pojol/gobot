const NodeTy = {
  // 行为树根节点
  Root: "RootNode",

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

  switch (ty) {
    case NodeTy.Root:
    case NodeTy.Loop:
    case NodeTy.Wait:
    case NodeTy.Selector:
    case NodeTy.Sequence:
      return false
    default:
      return true
  }

}

function IsActionNode(ty) {
  switch (ty) {
    case NodeTy.Root:
    case NodeTy.Loop:
    case NodeTy.Wait:
    case NodeTy.Selector:
    case NodeTy.Sequence:
    case NodeTy.Assert:
    case NodeTy.Condition:
      return false
    default:
      return true
  }
}


export { NodeTy, IsScriptNode, IsActionNode };
