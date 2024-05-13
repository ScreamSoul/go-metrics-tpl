package repositories

type RestoreMetricWrapper interface {
	MetricStorage
	Save()
	Load()
}
