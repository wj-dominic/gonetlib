module commandcenter

go 1.16

require (
	message v0.0.0
	ringbuffer v0.0.0
	session v0.0.0
	scv v0.0.0
)

replace (
	message v0.0.0 => ../message
	ringbuffer v0.0.0 => ../ringbuffer
	session v0.0.0 => ../session
	scv v0.0.0 => ../scv
)