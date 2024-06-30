package repositories

// RestoreMetricWrapper is an interface for a repository capable of saving and restoring dumps.
type RestoreMetricWrapper interface {
	MetricStorage
	Save()
	Load()
}
