package handler

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/hello", wrapper.requestParametersBinder)
	router.POST(baseURL+"/registration", wrapper.requestParametersBinder)
	router.POST(baseURL+"/login", wrapper.requestParametersBinder)
	router.GET(baseURL+"/profile", wrapper.requestParametersBinder)
	router.PATCH(baseURL+"/update-profile", wrapper.requestParametersBinder)
}
