package rwconwindows

import (
	"os"

	win "golang.org/x/sys/windows"
)

//Console type contains windows native fields, that can be accessed by its methods.
//Methods introduce a layer of abstraction over  windows native datatypes.
type Console struct {
	//pointer to underlying console :*os.File type
	filehandle *os.File
	//handle required inside library :win.Handle type
	winhandle win.Handle
	//info-field to store returns of libInternal functions :win.ConsoleScreenbufferInfo type
	info win.ConsoleScreenBufferInfo
}

//refreshes info field of registered Console;
func (c *Console) refreshConsoleScreenBufferInfo() {
	win.GetConsoleScreenBufferInfo(c.winhandle, &c.info)
}

//register console by handle (type: *os.File)
func (c *Console) Init(consolePtr *os.File) {
	c.filehandle = consolePtr
	c.winhandle = win.Handle(consolePtr.Fd())

	//Root.dimX, Root.dimY = c.WindowGetSize()
}

//returns filehandle of underlying Console
//
//Provides access to Functions of underlying console by getting its filehandle
func (c *Console) GetFH() *os.File {
	return c.filehandle
}

//returns info field of given Console
//
//reduces all types in info to primitives
func (c *Console) GetInfo() []int16 {
	c.refreshConsoleScreenBufferInfo()
	var (
		sizeX                   int16 = c.info.Size.X
		sizeY                   int16 = c.info.Size.Y
		posX                    int16 = c.info.CursorPosition.X
		posY                    int16 = c.info.CursorPosition.Y
		attr                    int16 = int16(c.info.Attributes)
		windowTop               int16 = c.info.Window.Top
		windowRight             int16 = c.info.Window.Right
		windowBottom            int16 = c.info.Window.Bottom
		windowLeft              int16 = c.info.Window.Left
		windowMaxSizeVertical   int16 = c.info.MaximumWindowSize.X
		windowMaxSizeHorizontal int16 = c.info.MaximumWindowSize.Y
	)
	return []int16{
		sizeX,
		sizeY,
		posX,
		posY,
		attr,
		windowTop,
		windowRight,
		windowBottom,
		windowLeft,
		windowMaxSizeVertical,
		windowMaxSizeHorizontal,
	}
}

//returns the current CursorPosition
//
//calls refreshConsoleScreenBufferInfo() then reads CursorPosition from info
func (c *Console) CursorGetPosition() (X, Y int16) {
	c.refreshConsoleScreenBufferInfo()
	X, Y = c.info.CursorPosition.X, c.info.CursorPosition.Y //extract X, Y from "CursorPosition" field in "info"
	return
}

//sets Cursorposition of registered Console to given position
func (c *Console) CursorSetPosition(posX, posY int16) {
	var posXY win.Coord = win.Coord{X: posX, Y: posY}
	win.SetConsoleCursorPosition(c.winhandle, posXY)
}

//returns horizontal and vertical size of the underlying consolewindow
func (c *Console) WindowGetSize() (horizontal, vertical int16) {
	horizontal = c.info.Window.Right
	vertical = c.info.Window.Bottom
	return
}

//draws given string into console at given position
//
//cReturn = true : Cursor is returned to initial position after write
func (c *Console) DrawAt(posX, posY int16, str string, cReturn bool) {
	if cReturn { //cursorReturn = true
		var x, y int16 = c.CursorGetPosition() //save initial CursorPosition
		c.CursorSetPosition(posX, posY)
		c.GetFH().WriteString(str)
		c.CursorSetPosition(x, y) //return Cursor
	} else { //cursorReturn = false
		c.CursorSetPosition(posX, posY)
		c.GetFH().WriteString(str)
	}
}

//moves Cursor by given offset
//
//calls CursorGetPosition to get posX, posY of Cursor
//calls SetConsoleCursorPosition to set new Position with offset
func (c *Console) CursorMove(offsetX, offsetY int16) {
	var posX, posY int16 = c.CursorGetPosition() //get posX & posY from ConsoleScreenBufferInfo by calling CursorGetPosition
	posX += offsetX                              //add offsetX
	posY += offsetY                              //addoffsetY
	if posX < 0 {                                //prevent negative final position
		posX = 0
	}
	if posY < 0 { //prevent negative final position
		posY = 0
	}
	var posXY win.Coord = win.Coord{X: posX, Y: posY}
	win.SetConsoleCursorPosition(c.winhandle, posXY) //move cursor to new position
}

/*** ELEMENT SYSTEM ***/

type Node struct {
	uid string
	X   int16
	Y   int16

	cNodes map[string]*Node
	cText  []Text
}
type Text struct {
	parentNode *Node
	offsetX    int16
	offsetY    int16
	value      string
}

//returns a New Nodeobject with given parameters
func NodeNew(uid string, posX, posY int16) Node {
	var newNode Node = Node{
		uid:    uid,
		X:      posX,
		Y:      posY,
		cNodes: map[string]*Node{},
		cText:  []Text{},
	}
	return newNode
}

//orders a nodeobject (caller) as the child of another nodeobject (argument)
//
//position of child is relative to parent-position
func (n *Node) NodeSetParent(parent *Node) {
	parent.cNodes[n.uid] = n
	n.X += parent.X
	n.Y += parent.Y
}

//Adds a Textfield at a relative position to a given node
func (n *Node) NodeNewText(value string, offsetX, offsetY int16) {
	n.cText = append(n.cText, Text{
		n,
		offsetX,
		offsetY,
		value,
	})
}

//draws all Textfields of a given Node
func (n *Node) draw() {
	/*
		for _, childText := range n.cText {
			fmt.Println("TextElement of '", n.uid, "' drawn at: {", n.X+childText.offsetX, "|", n.Y+childText.offsetY, "}")
		}
		fmt.Println("Node finished drawing:", n.uid)*/
}

//draws a node and all of its children
//
//calls draw() for a node and all of its children
func (n *Node) render() {
	n.draw()
	for _, childNode := range n.cNodes {
		childNode.render()
	}
}

//enables VTP for Console
//
//Virtual Terminal Processing has flawed cursorpositioning and needs further care when implemented
//can be used for color adjustments
func (c *Console) EnableVTProcessing() {
	var currentmode uint32
	win.GetConsoleMode(c.winhandle, &currentmode)
	win.SetConsoleMode(c.winhandle, currentmode|win.ENABLE_VIRTUAL_TERMINAL_PROCESSING|win.DISABLE_NEWLINE_AUTO_RETURN)
}
