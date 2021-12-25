

const Config = {
    driveAddr : "http://localhost:8888/",
    httpCode : `
local parm = {
    body = {},    -- request body
    timeout = "10s",
    headers = {},
}

local url = "http://127.0.0.1:port/api"
local http = require("http")

function execute()

    -- http post request
    res, errmsg = http.post(url, parm)
    if errmsg == nil then
        body = json.decode(res["body"])
        merge(meta, body.Body)
    end

end
`,
    assertCode : `
-- Write expression to return true or false
function execute()



end
`,
    conditionCode : `
-- Write expression to return true or false
function execute()



end
`
    ,
}


export default Config;