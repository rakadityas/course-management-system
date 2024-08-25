package config

// Server holds the server-related configuration.
type Server struct {
	HttpPort string
}

// Resource holds the resource-related configuration, such as databases.
type Resource struct {
	PrimaryDatabase string
}
