package openapi

type Err int32

const (
	Succ           Err = 200
	Fail           Err = 1000
	ErrContentRead     = 1000 + iota
	ErrWrongInput
	ErrJsonUnmarshal
	ErrJsonInvalid
	ErrPluginLoad
	ErrEnd
	ErrBreak
	ErrCantFindBot
	ErrCreateBot
	ErrEmptyBatch
	ErrTagsFormat
	ErrRunningErr
	ErrUploadConfig
	ErrGetConfig
)

var errmap map[Err]string = map[Err]string{
	ErrContentRead:   "failed to read request content",
	ErrJsonInvalid:   "wrong file format",
	ErrJsonUnmarshal: "json unmarshal err",
	ErrWrongInput:    "bad request parameter",
	ErrPluginLoad:    "failed to plugin load",
	ErrEnd:           "run to the end",
	ErrBreak:         "run to the break",
	ErrCantFindBot:   "can't find bot",
	ErrCreateBot:     "failed to create bot, the behavior tree file needs to be uploaded to the server before creation",
	ErrEmptyBatch:    "empty batch info",
}
