// Business Rules
package domain

// error messages:
// 1 is reserved to preserve error messages of external libraries
var errors = [...]string{
	1: "", // EXTERNAL_LIBRARY - used to handle external library errors
        2: "DBHOST environment variable required but not set",
        3: "DBPORT environment variable required but not set",
        4: "DBUSER environment variable required but not set",
        5: "DBPASS environment variable required but not set",
        6: "DBNAME environment variable required but not set",
	7: "DB connection is closed",
        8: "Rate Limit Exceed",
	9: "The requested field does not exist in this collection",
	10: "Item Id is empty - please provide a value",
	11: "Item Name is empty - please provide a value",
	12: "Item Quantity is empty - please provide a value",
}
