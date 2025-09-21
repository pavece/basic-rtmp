package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Registry = prometheus.NewRegistry()

	VideoIngress = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "video_ingress_bytes",
		Help: "Total amount of video bytes ingested by the server",
	})

	AudioIngress = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "audio_ingress_bytes",
		Help: "Total amount of audio bytes ingested by the server",
	})

	ActiveStreams = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "active_streams",
		Help: "Current number of active streams",
	})

	ObjectStoreUploads = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "object_store_uploads",
		Help: "Number of objects (resulting HLS chunks and lists) uploaded to object store",
	})
)

func init() {
	Registry.MustRegister(
		VideoIngress,
		AudioIngress,
		ActiveStreams,
		ObjectStoreUploads,
	)
}
