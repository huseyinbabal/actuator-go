package actuator

import "time"

type Request[T any] struct {
}

type Response[T any] struct {
}

type InfoResponse struct {
	Git struct {
		Branch string `json:"branch"`
		Commit struct {
			Id   string    `json:"id"`
			Time time.Time `json:"time"`
		} `json:"commit"`
	} `json:"git"`
	Build struct {
		Artifact string `json:"artifact"`
		Version  string `json:"version"`
		Group    string `json:"group"`
	} `json:"build"`
}

type HealthResponse struct {
	Status     string `json:"status"`
	Components struct {
		Broker struct {
			Status     string `json:"status"`
			Components struct {
				Us1 struct {
					Status  string `json:"status"`
					Details struct {
						Version string `json:"version"`
					} `json:"details"`
				} `json:"us1"`
				Us2 struct {
					Status  string `json:"status"`
					Details struct {
						Version string `json:"version"`
					} `json:"details"`
				} `json:"us2"`
			} `json:"components"`
		} `json:"broker"`
		Db struct {
			Status  string `json:"status"`
			Details struct {
				Database        string `json:"database"`
				ValidationQuery string `json:"validationQuery"`
			} `json:"details"`
		} `json:"db"`
		DiskSpace struct {
			Status  string `json:"status"`
			Details struct {
				Total     int64  `json:"total"`
				Free      int64  `json:"free"`
				Threshold int    `json:"threshold"`
				Path      string `json:"path"`
				Exists    bool   `json:"exists"`
			} `json:"details"`
		} `json:"diskSpace"`
	} `json:"components"`
}

type HeapDumpResponse struct {
	Location string `json:"location"`
}
