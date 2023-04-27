const NodeTy = {
  // 行为树根节点
  Root: "RootNode",

  // 脚本动作节点
  Action: "ActionNode",
  // 脚本条件节点
  Condition: "ConditionNode",

  // 循环次数节点
  Loop: "LoopNode",
  // 阻塞等待节点
  Wait: "WaitNode",

  Selector: "SelectorNode",
  Sequence: "SequenceNode",
  Parallel: "ParallelNode",
};

function IsScriptNode(ty: string) {

  switch (ty) {
    case NodeTy.Root:
    case NodeTy.Loop:
    case NodeTy.Wait:
    case NodeTy.Selector:
    case NodeTy.Sequence:
    case NodeTy.Parallel:
      return false
    default:
      return true
  }

}

function IsPresetNode(ty : string) {
  switch (ty) {
    case NodeTy.Root:
    case NodeTy.Loop:
    case NodeTy.Wait:
    case NodeTy.Selector:
    case NodeTy.Sequence:
    case NodeTy.Parallel:
    case NodeTy.Condition:
      return false
    default:
      return true
  }
}

function IsActionNode(ty: string) {
  switch (ty) {
    case NodeTy.Root:
    case NodeTy.Loop:
    case NodeTy.Wait:
    case NodeTy.Selector:
    case NodeTy.Sequence:
    case NodeTy.Condition:
    case NodeTy.Parallel:
      return false
    default:
      return true
  }
}


export { NodeTy, IsScriptNode, IsActionNode, IsPresetNode };  