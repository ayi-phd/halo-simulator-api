package domain

type EquipmentType string

const (
	EquipmentTypeFuelTruck     EquipmentType = "fuel_truck"
	EquipmentTypePushbackTug   EquipmentType = "pushback_tug"
	EquipmentTypeBeltLoader    EquipmentType = "belt_loader"
	EquipmentTypeCateringTruck EquipmentType = "catering_truck"
)

type EquipmentStatus string

const (
	EquipmentStatusAvailable EquipmentStatus = "available"
	EquipmentStatusBusy      EquipmentStatus = "busy"
	EquipmentStatusBroken    EquipmentStatus = "broken"
)

type GroundEquipment struct {
	ID       string          `json:"id"`
	Type     EquipmentType   `json:"type"`
	Status   EquipmentStatus `json:"status"`
	Location string          `json:"location"`
}
