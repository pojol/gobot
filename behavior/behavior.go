package behavior

type BehaviorType int

const (
	PostTy BehaviorType = iota + 1
	JumpTy
	DelayTy
)

/*
	{
		"behavior" : "POST",
		"url" : "",
		"name" : "",
		"script" : "",
		"param" : {

		}
	}
*/
type POST interface {
	Do() error
}

type TCP interface {
	Send() error
}

type Delay interface {
}

/*
	{
		"behavior" : "SELECT",
		"name" : "",
		"script" : "",
	}
*/
type Select interface {
}
