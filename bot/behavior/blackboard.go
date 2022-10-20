package behavior

type ThreadInfo struct {
	Num    int    `json:"num"`
	ErrMsg string `json:"errmsg"`
	CurNod string `json:"curnod"`
	Change string `json:"change"`
}

func (ti *ThreadInfo) reset() {
	ti.ErrMsg = ""
	ti.CurNod = ""
	ti.Change = ""
}

type Blackboard struct {
	Nods      []INod
	Threadlst []ThreadInfo

	end bool
}

func (b *Blackboard) GetOpenNods() []INod {
	return b.Nods
}

func (b *Blackboard) Append(nods []INod) {
	b.Nods = append(b.Nods, nods...)
}

func (b *Blackboard) Reset() {
	b.Nods = b.Nods[:0]
}

func (b *Blackboard) ThreadAdd(num int) {
	b.Threadlst = append(b.Threadlst, ThreadInfo{Num: num})
}

func (b *Blackboard) ThreadRmv(num int) {

}

func (b *Blackboard) ThreadFillInfo(info ThreadInfo) {
	for k, v := range b.Threadlst {
		if v.Num == info.Num {
			b.Threadlst[k] = info
		}
	}
}

func (b *Blackboard) ThreadInfoReset() {
	for k := range b.Threadlst {
		b.Threadlst[k].reset()
	}
}

func (b *Blackboard) ThreadInfo() []ThreadInfo {

	lst := []ThreadInfo{}

	for _, v := range b.Threadlst {
		if v.CurNod != "" {
			lst = append(lst, v)
		}
	}

	return lst
}

func (b *Blackboard) ThreadCurNum() int {
	num := 1

	for _, v := range b.Threadlst {
		if v.Num > num {
			num = v.Num
		}
	}

	return num
}

func (b *Blackboard) End() {
	b.end = true
}
