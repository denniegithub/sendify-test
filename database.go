package main

type DAO interface {
	CreateShipment(shipment *Shipment)
	ListShipmentsByCustomer(customerId string) []*Shipment
}

type MemoryDB struct {
	db map[string][]*Shipment
}

func initDB() *MemoryDB {
	db := make(map[string][]*Shipment)
	return &MemoryDB{db: db}
}

type Shipment struct {
	CustomerId string  `json:"customerId"`
	Sender     string  `json:"sender"`
	Receiver   string  `json:"receiver"`
	Weight     float32 `json:"weight"`
	Price      float32 `json:"price,omitempty"` // Read-only
}

func (m *MemoryDB) CreateShipment(shipment *Shipment) {

	records, ok := m.db[shipment.CustomerId]
	if !ok {
		m.db[shipment.CustomerId] = []*Shipment{shipment}
	} else {
		records = append(records, shipment)
		m.db[shipment.CustomerId] = records
	}

}

func (m *MemoryDB) ListShipmentsByCustomer(customerId string) []*Shipment {
	records, ok := m.db[customerId]
	if !ok {
		return []*Shipment{}
	}

	return records
}
