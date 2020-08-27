package yidoc

func getDescAndRequired(v ...interface{}) (string, bool) {
	var (
		desc     string
		required bool
	)
	if len(v) > 0 {
		desc, _ = v[0].(string)
	}
	if len(v) > 1 {
		required, _ = v[1].(bool)
	}
	return desc, required
}
