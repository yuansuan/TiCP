package syncpolicy

var DefaultRules = []*Rules{
	{
		Method:   RuleMethodPrefix,
		Express:  "d3plot",
		Compress: "lz",
	},
}

var DefaultPolicy *Policy

func init() {
	DefaultPolicy, _ = NewPolicy(DefaultRules)
}
