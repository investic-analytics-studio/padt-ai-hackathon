package constant

import "fmt"

// GetTimeFilter returns a SQL time filter based on the given time range
func GetTimeFilter(timeRange string) string {
	switch timeRange {
	case "1d":
		return " AND date >= NOW() - INTERVAL 1 DAY"
	case "7d":
		return " AND date >= NOW() - INTERVAL 7 DAY"
	case "1m":
		return " AND date >= NOW() - INTERVAL 1 MONTH"
	case "3m":
		return " AND date >= NOW() - INTERVAL 3 MONTH"
	case "6m":
		return " AND date >= NOW() - INTERVAL 6 MONTH"
	case "1year":
		return " AND date >= NOW() - INTERVAL 1 YEAR"
	case "2year":
		return " AND date >= NOW() - INTERVAL 2 YEAR"
	case "all":
		return " "
	default:
		return fmt.Sprintf(" AND date >= NOW() - INTERVAL %s", timeRange)
	}
}
