package container

// Init exists to allow container init() functions to run when
// invoked from the capabilities Init function, which is
// itself invoked by the Lumogon command handler.
func Init() {
}
