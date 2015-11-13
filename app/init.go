package app

import (
	"github.com/revel/revel"
	"net"
	"os"
	"encoding/json"
	"fmt"
	"os/user"
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
