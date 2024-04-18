
local parm = {
  body = {
      Token = meta.Token,
      HeroID = "joy"
  },    -- request body
  timeout = "10s",
  headers = {},
}

local url = REMOTE .. "/base/hero.lvup"
local http = require("http")

function execute()

  -- http post request
  res, errmsg = http.post(url, parm)
  if errmsg == nil then
    body = json.decode(res["body"])
    merge(meta, body.Body)
  end

  return state.Succ, body.Body
end
