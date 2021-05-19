package nomad

const (
	lokiDriverName    = "loki"
	elasticDriverName = "elastic/elastic-logging-plugin:7.12.1"
	indexName         = "scarlet-fairy-workloads"
)

func loggingConfig(url string, id string) map[string]interface{} {
	return map[string]interface{}{
		"type": elasticDriverName,
		"config": []map[string]string{
			{
				"hosts": url,
				"name":  id,
				"index": indexName,
			},
		},
	}
}
