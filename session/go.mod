module session

go 1.16

require (
	message v0.0.0
	ringbuffer v0.0.0
)

replace (
	message v0.0.0 => ../message
	ringbuffer v0.0.0 => ../ringbuffer
)