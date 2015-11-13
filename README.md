# telerobot
Telepresence robot server prototype written in Go and running on Revel. To be used in conjunction with a differential wheeled robot powered by a Particle Photon (or Core) microcontroller. 

## Configuration
Copy `telerobot_conf.json.sample` to `~/.telerobot_conf.json` and update with the device id and access token of the particle microcontroller.

### Start the web server:

    revel run github.com/kldavis4/telerobot

   Run with <tt>--help</tt> for options.

### Go to http://localhost:9000/

This page allows control via the virtual joystick. There are two status indicators at the top left. The top indicator is the device status (green = connected, red = not connected). The bottom indicator shows the motion command server status (yellow = running).

### Got to http://localhost:9000/program

This page allows control via a list of motion commands.

