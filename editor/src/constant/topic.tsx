

const Topic = {
    NodeRmv : "topic_node_remove",
    NodeHistoryClean : "topic_history_clean",   // 清空 cmd list
    NodeUpdateParm: "topic_node_update_parm",

    NodeGraphClick: "topic_node_graph_click",  // 在editor的编辑器视图中点击节点
    NodeEditorClick: "topic_node_editor_click",    //在editor的edit框中点击节点类型

    ThemeChange: "theme_change",

    BotsUpdate: "topic_bots_update",   // 机器人列表需要更新

    PrefabUpdateAll: "topic_prefab_update_all",

    FileSave : "topic_file_save",
    FileLoadDraw: "topic_file_load_draw",
    FileLoadRedraw : "topic_file_load_graph",

    ReportSelect : "topic_report_select",

    // c2s
    DebugUpload : "topic_debug_upload",    // 上传行为树模版文件
    DebugCreate : "topic_debug_create",    // 基于某个模版，创建一个可运行的 Bot
    DebugStep : "topic_debug_step",        // 单步运行某个 Bot
    
    // s2c
    DebugUpdateBlackboard : "topic_debug_update_blackboard",    // 将运行 Bot 返回的数据写入到 blackboard
    DebugUpdateChange : "topic_debug_update_change",
    DebugFocus: "topic_debug_focus", // 加亮显示当前运行到的某个 Node

}

export default Topic;