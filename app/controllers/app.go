package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"math"
	"strings"
	"time"
	"github.com/kldavis4/telerobot/app"
)

//Basic response structure for API requests
type ApiResponse struct {
	Code int
	Message string
	Body string
}

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Program() revel.Result {
	return c.Render()
}

// Executes a set of motion commands specified one per line
// Each command consists of two 3 digit numeric value with a max value of 255 prefixed with a plus or minus symbol
// and a duration specified in milliseconds
// Only one program can execute at a time
func (c App) ExecuteProgram(commandset string) revel.Result {
	if app.ProgramExecuting == false {
		commands := strings.Split(commandset,"\n")

		go executeCommands(commands)

		return c.RenderJson(ApiResponse{200,"Success",fmt.Sprintf("%d instructions sent",len(commands))})
	} else {
		return c.RenderJson(ApiResponse{400,"Failure","Program execution in progress"})
	}
}

//Updates the current move command to the one specified
//dx, dy are collected from a joystick.js virtual joystick
func (c App) Move(dx int, dy int) revel.Result {
	if app.ProgramExecuting == true {
		return c.RenderJson(ApiResponse{400,"Failure","Program executing"})
	}
	
	state := formatJoystickMotion(dx,dy)

	if strings.Compare(state,app.State) != 0 {
		app.Dirty = true
		app.State = state
	}

	return c.RenderJson(ApiResponse{200,"Success",app.State})
}

//Queries the current device status via the particle api
func (c App) Status() revel.Result {
	resp, err := http.Get(fmt.Sprintf("https://api.particle.io/v1/devices/%s?access_token=%s",app.Config.DeviceId,app.Config.AccessToken))

	if ( err == nil ) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		
		if ( err == nil ) {
			var f map[string]interface{}
			err := json.Unmarshal(body, &f)

			if ( err == nil ) {
				if f["connected"] == true {
					return c.RenderJson(ApiResponse{200,"Success","connected"})
				} else {
					return c.RenderJson(ApiResponse{200,"Success","not_connected"})
				}
			} else {
				fmt.Printf(fmt.Sprintf("Error parsing json: %s", err))
			}
			return c.RenderJson(ApiResponse{200,"Success",string(body[:])})
		} else {
			fmt.Printf(fmt.Sprintf("Error reading response body: %s\n", err));
			return c.RenderJson(ApiResponse{500,"Failure",fmt.Sprintf("Error reading response body: %s", err)})
		}
	} else {
		fmt.Printf(fmt.Sprintf("Error during request: %s\n", err));
		return c.RenderJson(ApiResponse{500,"Failure",fmt.Sprintf("Error during request: %s", err)})	
	}	
}

//Parse the joystick.js dx / dy into a format that that robot can use
func formatJoystickMotion(dx int, dy int) string {
	rightMagnitude := -float64(dy)
	leftMagnitude := -float64(dy)

	//Set a deadzone between -15 and +15 to make driving straight easier
	if dx > 15 {
		leftMagnitude = leftMagnitude * (1.0 - float64(dx) / 255.0)
	} else if dx < -15 {
		rightMagnitude = rightMagnitude * (1.0 + float64(dx) / 255.0)
	}
	
	return fmt.Sprintf("%0+4d%0+4d",int(math.Floor(leftMagnitude)),int(math.Floor(rightMagnitude)))
}

//Execute a command set
func executeCommands(commands []string) {
	fmt.Println(fmt.Sprintf("Executing %d commands", len(commands)))
	app.ProgramExecuting = true
	for _, command := range commands {
		parts := strings.Split(command," ")
		
		d, err := time.ParseDuration(fmt.Sprintf("%sms", parts[2]))
		if err != nil {
			fmt.Println("Invalid duration ", parts[2])
			continue
		}
		time.Sleep(d)
		
		app.State = fmt.Sprintf("%s%s",parts[0],parts[1])
		app.Dirty = true
		fmt.Println(app.State)
	}
	app.ProgramExecuting = false
}
