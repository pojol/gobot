package behavior

type BehaviorType int

const (
	PostTy BehaviorType = iota + 1
	JumpTy
	DelayTy
)

/*
	tree compose
	// 通过树编排文件执行逻辑 （脚本文件调度
	// 数据文件
		// 有 tmp-time-files - 通过数据文件，构建脚本文件；然后执行
		// 无 script-files - 直接使用现有脚本

	{
		"root":"",
		"name":"",
		"script":"ONCE | LOOP",

		"children" : [

			{
				"behavior" : "SELECT",
				"name" : "",
				"script" : " if meta.Token == "" return 0 else return 1,2 ",
				"children" : [
					{
						"behavior" : "POST",
						"name" : "login",
						"script" : "登录"
					},
					{
						"behavior" : "POST",
						"name" : "getinfo",
						"script" : "获取账户信息"
					}
				]
			}

		]
	}
*/

/*
	{
		"behavior" : "POST",
		"name" : "",
		"script" : "",
	}
*/
type IPOST interface {
	Do([]byte) ([]byte, error)
}

type ISend interface {
	Do() error
}

/*
	{
		"behavior" : "DELAY",
		"dura" : 100, // ms
	}
*/
type IDelay interface {
	Do() error
}

/*
	{
		"behavior" : "ASSERT",
		"name" : "",
		"script" : "xxx.lua",
	}
*/
type IAssert interface {
	Do() error
}

/*
	{
		"behavior" : "SELECT",
		"name" : "",
		"script" : "xxx.lua",
	}

*/
type Select interface {
}
