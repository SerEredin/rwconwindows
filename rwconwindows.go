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

//enables VTP for Console
//
//Virtual Terminal Processing has flawed cursorpositioning and needs further care when implemented
//can be used for color adjustments
func (c *Console) EnableVTProcessing() {
	var currentmode uint32
	win.GetConsoleMode(c.winhandle, &currentmode)
	win.SetConsoleMode(c.winhandle, currentmode|win.ENABLE_VIRTUAL_TERMINAL_PROCESSING|win.DISABLE_NEWLINE_AUTO_RETURN)
}

/*********************elements***********************
type element struct {
	posX int16
	posY int16

	parent   *element
	children []*element

	contents string
}

//DocRoot is the RootElement
//
//new elements not added to a certain parent are added to DocRoot
var DocRoot element = element{
	posX: 0,
	posY: 0,
}

func (p *element) NewElement(posX, posY int16, val string) (newEl element) {
	newEl = element{ //							create new element
		posX:     posX, //at given X relative to parent
		posY:     posY, //at given Y relative to parent
		parent:   p,    //as a child of the parent NewElement was called on
		contents: val,  //with the desired value
	}
	p.children = append(p.children, &newEl)
	return
}
func (p *element) DrawElement(c Console, posX, posY int16) {
	for i := 0; i < len(p.children); i++ {
		var child *element = p.children[i] //get ref to child-element

		child.DrawElements(c, p.posX+child.posX, p.posY+child.posY)
	}

}
*/
