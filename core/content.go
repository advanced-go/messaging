package core

// Map - slice of any content to be included in a message
type Map map[string][]any

// Credentials - type for a credentials function
type Credentials func() (username string, password string, err error)

// Resource - struct for a resource
type Resource struct {
	Uri string
}

// AccessCredentials - access function for Credentials in a message
func AccessCredentials(msg *Message) Credentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(Credentials); ok {
			return fn
		}
	}
	return nil
}

// AccessResource - access function for a resource in a message
func AccessResource(msg *Message) Resource {
	if msg == nil || msg.Content == nil {
		return Resource{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(Resource); ok {
			return url
		}
	}
	return Resource{}
}
