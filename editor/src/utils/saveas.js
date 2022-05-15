

export function SaveAs(blob, filename) {

    if (window.navigator.msSaveOrOpenBlob) {
        navigator.msSaveBlob(blob, filename)
    } else {

        // 上面这个是创建一个blob的对象连链接，
        // 创建一个链接元素，是属于 a 标签的链接元素，所以括号里才是a，
        var link = document.createElement("a");
        var body = document.querySelector("body")

        link.href = window.URL.createObjectURL(blob);
        link.download = filename
        
        // firefox
        link.style.display = "node"
        body.appendChild(link)
      
        // 使用js点击这个链接
        link.click();
        body.removeChild(link)

        window.URL.revokeObjectURL(link.href)
    }
    
}