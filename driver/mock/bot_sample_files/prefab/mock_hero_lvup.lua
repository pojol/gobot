
local parm = {
  body = {
      Token = bot.Token,
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
    merge(bot, body.Body)
  end

  return state.Succ, body.Body
end
