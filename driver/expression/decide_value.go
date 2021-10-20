package expression

func (e *Expression) decideWithValue(values ...interface{}) bool {
	var b bool

	switch e.Symbol {
	case "$and":
	case "$or":
	case "$eq":
	case "$ne":
	case "$gt":
	default:
		println("decide", e.Symbol)
	}

	return b
}

/*
	$eq.Token
	$eq.Market.ID
	$and.$ne.Token.$gt.Gold

	map[id] = "$eq.Token"

	switch symbol {
	case $eq.Token:
	}
*/

func (eg *ExpressionGroup) DecideWithValue(values ...interface{}) bool {
	return eg.Root.decideWithValue(values...)
}
