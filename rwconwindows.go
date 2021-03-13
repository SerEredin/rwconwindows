package rwconwindows

import (
	"os"

	win "golang.org/x/sys/windows"
)

//Console type contains windows native fields, that can be accessed by its methods.
//Methods introduce a layer of abstraction over  windows native datatypes.
type Console struct {
	filehandle *os.File
	winhandle  win.Handle
	info       win.ConsoleScreenBufferInfo
}

//returns all primitives of info as int16
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

//returns filehandle of Console (*os.file)
func (c *Console) GetFH() *os.File {
	return c.filehandle
}

//returns height
//
//Rows from top to bottom of current window. Resizing window may change the actual value
func (c *Console) GetHeight() (bottom int16) {
	c.refreshConsoleScreenBufferInfo() //refresh c.info
	bottom = c.info.Window.Bottom      //row count to bottom
	return
}

//returns width
//
//Columns from left side to right side
func (c *Console) GetWidth() (right int16) {
	c.refreshConsoleScreenBufferInfo() //refresh c.info
	right = c.info.Window.Right        //column count to right
	return
}

//returns the current CursorPosition
//
//calls refreshConsoleScreenBufferInfo to get current posXY from c.info
func (c *Console) GetCursorPosition() (X, Y int16) {
	c.refreshConsoleScreenBufferInfo()                      //refresh c.info
	X, Y = c.info.CursorPosition.X, c.info.CursorPosition.Y //extract X, Y from "CursorPosition" field in "info"
	return
}

//register console by handle (type: *os.File)
func (c *Console) Init(consolePtr *os.File) {
	c.filehandle = consolePtr
	c.winhandle = win.Handle(consolePtr.Fd())
}

//refreshes c.info field of registered Console;
func (c *Console) refreshConsoleScreenBufferInfo() {
	win.GetConsoleScreenBufferInfo(c.winhandle, &c.info)
}

//draws given string into console at given position
func (c *Console) DrawAt(posX, posY int16, str string) {
	c.CursorSetPosition(posX, posY)
	c.GetFH().WriteString(str)
}

//sets Cursorposition of registered Console to given position
func (c *Console) CursorSetPosition(posX, posY int16) {
	var posXY win.Coord = win.Coord{X: posX, Y: posY}
	win.SetConsoleCursorPosition(c.winhandle, posXY)
}

//moves Cursor of registered Console by given offset
//
//calls GetCursorPosition to get posX, posY of Cursor
//then sets new position using SetConsoleCursorPosition
func (c *Console) CursorMove(offsetX, offsetY int16) {
	var posX, posY int16 = c.GetCursorPosition() //get posX & posY from ConsoleScreenBufferInfo by calling GetCursorPosition
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
