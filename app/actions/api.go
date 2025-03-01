// app/actions/api.go
package actions

// DataItem represents a data structure from API
type DataItem struct {
	Title       string
	Description string
}

// FetchData simulates an API call to get data
func FetchData() []DataItem {
	// In a real app, this would call an API or database
	return []DataItem{
		{Title: "API Item 1", Description: "Description from API call"},
		{Title: "API Item 2", Description: "Another description from API"},
	}
}
