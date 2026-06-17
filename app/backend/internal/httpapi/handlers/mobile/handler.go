// internal/httpapi/handlers/mobile/handler.go
package mobile

type Handler struct {
	PairingService  PairingService
	LocationService LocationService
	DeviceService   DeviceService
}
