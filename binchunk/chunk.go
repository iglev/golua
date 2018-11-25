package binchunk

const (
	// LuaSignature lua binchunk header signature
	LuaSignature = "\x1bLua"
	// LuacVersion luac version (5.3 -> 5*16+3 -> 0x53)
	LuacVersion = 0x53
	// LuacFormat luac format
	LuacFormat = 0
	// LuacData luac data (0x19 93 0D 0A 1A 0A)
	LuacData = "\x19\x93\r\n\x1a\n"
	// CIntSize c int size
	CIntSize = 4
	// CSizetSize c size_t size
	CSizetSize = 8
	// InstructionSize instruction size
	InstructionSize = 4
	// LuaIntSize lua int size
	LuaIntSize = 8
	// LuaNumberSize lua number size
	LuaNumberSize = 8
	// LuacInt luac int (check big/little endian)
	LuacInt = 0x5678
	// LuacNum luac number (check double IEEE754)
	LuacNum = 370.5
)

const (
	// TagNil nil
	TagNil = 0x00
	// TagBoolean boolean
	TagBoolean = 0x01
	// TagNumber number
	TagNumber = 0x02
	// TagInteger integer
	TagInteger = 0x03
	// TagShortStr short str
	TagShortStr = 0x04
	// TagLongStr long str
	TagLongStr = 0x05
)

type binaryChunk struct {
	header
	sizeUpvalues byte
	mainFunc     *ProtoType
}

type header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntSize      byte
	luaNumberSize   byte
	luacInt         int64
	luacNum         float64
}

// ProtoType prot type
type ProtoType struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*ProtoType
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

// Upvalue upvalue
type Upvalue struct {
	Instack byte
	Idx     byte
}

// LocVar local var
type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

/*
// Undump ...
func Undump(data []byte) *ProtoType {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte()
	return reader.readPorot("")
}
*/
