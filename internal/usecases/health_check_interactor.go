package usecases

type HealthCheckInteractor interface {
	PerformHealthCheck() HealthStatus
}

type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type healthCheckInteractorImpl struct{}

func NewHealthCheckInteractor() HealthCheckInteractor {
	return &healthCheckInteractorImpl{}
}

func (h *healthCheckInteractorImpl) PerformHealthCheck() HealthStatus {
	return HealthStatus{
		Status:  "Healthy",
		Message: "Success",
	}
}
