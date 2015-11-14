# telerobot
Telepresence robot server prototype written in Go and running on Revel. To be used in conjunction with a differential wheeled robot powered by a Particle Photon (or Core) microcontroller. 

## Configuration
Copy `telerobot_conf.json.sample` to `~/.telerobot_conf.json` and update with the device id and access token of the particle microcontroller.

### Start the web server:

    revel run github.com/kldavis4/telerobot

   Run with <tt>--help</tt> for options.

### Go to http://localhost:9000/

This page allows control via the virtual joystick. There is a device status indicator at the top left (green = connected, red = not connected, purple = web application error).

### Got to http://localhost:9000/program

This page allows control via a list of motion commands.

