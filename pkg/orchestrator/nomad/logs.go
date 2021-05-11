package nomad

const (
	lokiDriverName = "loki"
)

func loggingConfig(lokiUrl string) map[string]interface{} {
	return map[string]interface{}{
		"type": lokiDriverName,
		"config": []map[string]string{
			{
				"loki-url": lokiUrl,
			},
		},
	}
}
