package state

import (
	"github.com/iglev/glua/api"
	"github.com/iglev/glua/binchunk"
	"github.com/iglev/glua/vm"
)

// Load - lua_load
func (l *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk) // todo
	c := newLuaClosure(proto)
	l.stack.push(c)
	return 0
}

// Call - lua_call
func (l *luaState) Call(nArgs, nResults int) {
	val := l.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		if c.proto != nil {
			l.callLuaClosure(nArgs, nResults, c)
		} else {
			l.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not function")
	}
}

func (l *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+api.LUA_MINSTACK, l)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := l.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	l.stack.pop()

	// run closure
	l.pushLuaStack(newStack)
	r := c.goFunc(l)
	l.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(r)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
}

func (l *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+api.LUA_MINSTACK, l)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := l.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	l.pushLuaStack(newStack)
	l.runLuaClosure()
	l.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
}

func (l *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(l.Fetch())
		inst.Execute(l)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}