package inventory

// TODO: the "Shape" can be one of 5000+ shapes. See the iOS app "Steel".
// We are also allowed to import the "Steel" DB to this app

// Inventory is an inventory item
type Inventory struct {
	ID        string
	ProjectID string
	Stage     Stage
	Size      int32
	Length    int32
	Grade     int32
	Shape     string
	Passed    bool
	Sequence  int32
	Priority  int32
}

// Stage is what stage the inventory item is in
type Stage struct {
	Raw       bool
	InProcess bool
	OnHold    bool
	Finished  bool
}
