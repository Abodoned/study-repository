package repositor

type EquipmentType struct {
	Name string
	Cost Coal
}

var EquipmentCatalog = map[string]EquipmentType{
	"pickaxe": {
		Name: "pickaxe",
		Cost: 3000,
	},
	"ventilation": {
		Name: "ventilation",
		Cost: 15000,
	},
	"car": {
		Name: "car",
		Cost: 50000,
	},
}

func GetEquipmentType(name string) (EquipmentType, error) {
	item, ok := EquipmentCatalog[name]
	if !ok {
		return EquipmentType{}, ErrorEquipmentNotFound
	}
	return item, nil
}
