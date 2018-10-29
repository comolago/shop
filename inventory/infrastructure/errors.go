package infrastructure

var errors = [...]string{
	1: "", // EXTERNAL_LIBRARY - used to handle external library errors
	2: "DBHOST environment variable required but not set",
	3: "DBPORT environment variable required but not set",
	4: "DBUSER environment variable required but not set",
	5: "DBPASS environment variable required but not set",
	6: "DBNAME environment variable required but not set",
}
