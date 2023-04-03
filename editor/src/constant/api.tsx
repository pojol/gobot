const Api = {
    //
    FileBlobUpload: "file.uploadBlob",
    FileTxtUpload: "file.uploadTxt",
    FileGet: "file.get",
    FileList: "file.list",
    FileRemove: "file.remove",
    FileSetTags: "file.setTags",

    /*
        Debug
    */
    DebugCreate: "debug.create",
    DebugStep: "debug.step",
    DebugInfo : "debug.info",

    /*
        Bot
    */
    BotRun: "bot.run",
    BotCreateBatch: "bot.batch",
    BotList: "bot.list",

    /*
        Prefab
    */
    PrefabUpload: "prefab.upload",
    PrefabList : "prefab.list",
    PrefabGet: "prefab.get",
    PrefabRemove : "prefab.rmv",
    PrefabSetTags : "prefab.setTags",

    /*
        Report
    */
    ReportInfo: "report.get",

    /*
        Config
    */
    ConfigSystemInfo : "config.sys.info",
    ConfigSystemSet : "config.sys.set",
    ConfigGlobalInfo : "config.global.info",
    ConfigGlobalSet : "config.global.set"
}


export default Api;