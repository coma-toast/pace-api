package inspection

import (
	"time"

	"github.com/coma-toast/pace-api/pkg/inventory"
)

// Inspection is an inspection report
type Inspection struct {
	InspectionID   string
	ProjectID      string
	Username       string
	StartTime      time.Time
	EndTime        time.Time
	InspectedParts []inventory.Inventory
}
