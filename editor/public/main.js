const { app, BrowserWindow, Menu } = require('electron');
const isDev = require('electron-is-dev');


// window对象的全局引用
let mainWindow
function createWindow() {

    Menu.setApplicationMenu(null)
    mainWindow = new BrowserWindow({ width: 1680, height: 1050 ,title:"Gobot"})

    // 开发环境
    //mainWindow.loadURL('http://localhost:3000/');

    // 生产环境 
    mainWindow.loadURL(isDev ? 'http://localhost:8000' : `file://${__dirname}/../dist/index.html`);

    // 打开开发者工具，默认不打开
    isDev && mainWindow.webContents.openDevTools();

    // 关闭window时触发下列事件.
    mainWindow.on('closed', function () {
        mainWindow = null
    })
}

app.on('ready', createWindow);

// 所有窗口关闭时退出应用.
app.on('window-all-closed', function () {
    if (process.platform !== 'darwin') {
        app.quit()
    }
})

app.on('activate', function () {

    if (mainWindow === null) {
        createWindow()
    }
})