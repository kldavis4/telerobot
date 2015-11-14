package app

import (
	"github.com/revel/revel"
	"net"
	"os"
	"encoding/json"
	"fmt"
	"os/user"
	"strings"
	"time"
)

//App configuration
type Configuration struct {
    DeviceId string			//Particle device id
	AccessToken string		//particle device access token
	MotionServerPort string	//Port where motion command server listens
}

var State string			//Current Move Command State
var Dirty bool				//Dirty flag for Move Command State
var Listener net.Listener	//Move Command Server Listener
var ProgramExecuting bool	//Flag to indicate if there is a program currently executing
var Config Configuration	//App configuration

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
	
	State = ""
	Dirty = true
	
	usr, err := user.Current()
    if err != nil {
        fmt.Println("error:", err)
    }
	
	file, err := os.Open(fmt.Sprintf("%s/.telerobot_conf.json", usr.HomeDir))
	if err != nil {
		fmt.Println("error:", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	revel.OnAppStart(initMotionServer)
}

func initMotionServer() {
	go startMotionServer()
}

//Start the motion command server
func startMotionServer() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s",Config.MotionServerPort))
	if err != nil {
	    fmt.Println("Error creating listener:", err.Error())
	    return
	}

	fmt.Println(fmt.Sprintf("Motion command server listening on %s",Config.MotionServerPort))

	Listener = ln

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
	        fmt.Println("Error accepting: ", err.Error())
	        return
		}
		go handleConnection(conn)
	}
}

// Handle connection requests to the server
func handleConnection(conn net.Conn) {
	buf := make([]byte, 128)

	for {
		if Dirty {	//Only send updates
			conn.Write([]byte(State))

			reqLen, err := conn.Read(buf)

			if err != nil {
				fmt.Println("Error reading:", err.Error())
			} else {
				ackStr := string(buf[0:reqLen])
				if strings.Contains(ackStr, "ACK") {
					Dirty = false
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Close the connection when you're done with it.
	conn.Close()
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
