package config

// FocusableItem defines an interface for items that can receive focus
type FocusableItem interface {
	Focus()
	Blur()
	IsFocused() bool
}

// NavigationManager handles focus navigation between items
type NavigationManager struct {
	Items        []FocusableItem
	CurrentFocus int
}

// NewNavigationManager creates a new navigation manager
func NewNavigationManager() *NavigationManager {
	return &NavigationManager{
		Items:        make([]FocusableItem, 0),
		CurrentFocus: -1,
	}
}

// AddItem adds a focusable item to the manager
func (nm *NavigationManager) AddItem(item FocusableItem) {
	nm.Items = append(nm.Items, item)
}

// Next focuses the next item in the list
func (nm *NavigationManager) Next() {
	if len(nm.Items) == 0 {
		return
	}

	// Blur current item
	if nm.CurrentFocus >= 0 && nm.CurrentFocus < len(nm.Items) {
		nm.Items[nm.CurrentFocus].Blur()
	}

	// Move to next item
	nm.CurrentFocus = (nm.CurrentFocus + 1) % len(nm.Items)
	nm.Items[nm.CurrentFocus].Focus()
}

// Previous focuses the previous item in the list
func (nm *NavigationManager) Previous() {
	if len(nm.Items) == 0 {
		return
	}

	// Blur current item
	if nm.CurrentFocus >= 0 && nm.CurrentFocus < len(nm.Items) {
		nm.Items[nm.CurrentFocus].Blur()
	}

	// Move to previous item
	if nm.CurrentFocus <= 0 {
		nm.CurrentFocus = len(nm.Items) - 1
	} else {
		nm.CurrentFocus--
	}
	nm.Items[nm.CurrentFocus].Focus()
}
