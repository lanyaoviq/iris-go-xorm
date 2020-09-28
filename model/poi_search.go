package model

type PoiSearch struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Latitude float64 `json:"latitude"`
	Geohash string `json:"geohash"`
	
}
