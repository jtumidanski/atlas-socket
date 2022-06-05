package request

// RequestReader -
type RequestReader struct {
	pos    int
	packet *Request
	Time   int64
}

// NewRequestReader -
func NewRequestReader(p *Request, time int64) RequestReader {
	return RequestReader{pos: 0, packet: p, Time: time}
}

func (r RequestReader) String() string {
	return r.packet.String()
}

// GetBuffer -
func (r *RequestReader) GetBuffer() []byte {
	return *r.packet
}

func (r *RequestReader) GetRestAsBytes() []byte {
	return (*r.packet)[r.pos:]
}

func (r *RequestReader) Skip(amount int) {
	if len(*r.packet)-(r.pos+amount) >= 0 {
		r.pos += amount
	}
}

// ReadByte -
func (r *RequestReader) ReadByte() byte {
	if len(*r.packet)-r.pos > 0 {
		return r.packet.readByte(&r.pos)
	}

	return 0
}

// ReadInt8 -
func (r *RequestReader) ReadInt8() int8 {
	if len(*r.packet)-r.pos > 0 {
		return r.packet.readInt8(&r.pos)
	}

	return 0
}

// ReadBool -
func (r *RequestReader) ReadBool() bool {
	if len(*r.packet)-r.pos > 0 {
		return r.packet.readBool(&r.pos)
	}

	return false
}

// ReadBytes -
func (r *RequestReader) ReadBytes(size int) []byte {
	if len(*r.packet)-r.pos >= size {
		return r.packet.readBytes(&r.pos, size)
	}

	return []byte{0}
}

// ReadInt16 -
func (r *RequestReader) ReadInt16() int16 {
	if len(*r.packet)-r.pos > 1 {
		return r.packet.readInt16(&r.pos)
	}

	return 0
}

// ReadInt32 -
func (r *RequestReader) ReadInt32() int32 {
	if len(*r.packet)-r.pos > 3 {
		return r.packet.readInt32(&r.pos)
	}

	return 0
}

// ReadInt64 -
func (r *RequestReader) ReadInt64() int64 {
	if len(*r.packet)-r.pos > 7 {
		return r.packet.readInt64(&r.pos)
	}

	return 0
}

// ReadUint16 -
func (r *RequestReader) ReadUint16() uint16 {
	if len(*r.packet)-r.pos > 1 {
		return r.packet.readUint16(&r.pos)
	}

	return 0
}

// ReadUint32 -
func (r *RequestReader) ReadUint32() uint32 {
	if len(*r.packet)-r.pos > 3 {
		return r.packet.readUint32(&r.pos)
	}

	return 0
}

// ReadUint64 -
func (r *RequestReader) ReadUint64() uint64 {
	if len(*r.packet)-r.pos > 7 {
		return r.packet.readUint64(&r.pos)
	}

	return 0
}

// ReadString -
func (r *RequestReader) ReadString(size int16) string {
	if len(*r.packet)-r.pos >= int(size) {
		return r.packet.readString(&r.pos, int(size))
	}

	return ""
}

func (r *RequestReader) ReadAsciiString() string {
	am := r.ReadInt16()
	return r.ReadString(am)
}

func (r *RequestReader) Position() int {
	return r.pos
}

func (r *RequestReader) Seek(offset int) {
	r.pos = offset
}

func (r *RequestReader) Available() int {
	return r.packet.Size() - r.pos
}
