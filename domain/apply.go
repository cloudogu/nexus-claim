package domain

// ApplyPlan executes a plan against the nexus api
func ApplyPlan(writer NexusAPIWriter, plan *Plan) error {
	return plan.Execute(writer)
}
