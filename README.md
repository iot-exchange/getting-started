# IoT Exchange - Getting Started

Welcome to the IoT Exchange Getting Started repository! 
This guide will help you quickly getting started with testing your device.

## Prerequisites

- An IoT Exchange account ([Sign up here](https://portal.iot-exchange.io))
- git / curl / bash installed to use the scripts provided in this repo

## Steps to Get Started

### 1. Sign Up or Log In
If you don't already have an account, sign up at [IoT Exchange](https://portal.iot-exchange.io). If you already have an account, simply log in.

### 2. Create or Select a Project
Once logged in:
- Navigate to the **Projects** section.
- Create a new project or select an existing one.

### 3. Add a Device
To register a new device:
- Go to the **Devices** section.
- Click **Add Device**.
- Set the name of the device to something you can easily identify.
- Set **Network** to `Internet` and **Protocol** to `http`.
- Finish adding the device.

### 4. Generate a Certificate
After adding your device:
- Click **Generate Certificate** for your device and copy the generated certificate.
- Open a text editor of your choice and paste the generated certificate.
- Save the file as `\<your device\>.pem`.

### 5. Add an Integration
Integrations are external applications that will receive the messages from
all devices in the project.
For testing we are going to add a test integration with the free service
[webhook.site](https://webhook.site) that allows you to receive webhooks in
your browser in realtime.
- Go to https://webhook.site and click on "Your unique URL" to copy the unique url it has generated for you.
- Go back to the IoT Exhchange and navigate to the integrations of your project.
- Click **Add Integration**.
- Paste your unique webhook url into the URI field. You can leave the rest as-is.
- Finish adding the Integration by clicking on **Save**

### 6. Send Data Using the Provided Script
This repository contains a script to send data from your device to IoT Exchange.
It uses `curl` to send uplink messages to the iot-exchange as if it was
coming from an actual IoT device.
This uses the `http` protocol; for actual IoT devices you want probably want
to use `coap` or `udp` with dtls as that uses way less data and is more
suited for IoT networks.

#### Steps to run the script:
1. Clone this repository:
   ```sh
   git clone https://github.com/iot-exchange/getting-started.git
   cd getting-started
   ```
2. Run the script with :
   ```sh
   echo -n "\<Your message goes here\>" | ./uplink.sh `/path/to/\<your device\>.pem`
   ```

### 7. View your message
- Go to https://webhook.site this should now display your message
- It should look something like this: ```json
{
  "timestamp": "1741512267750437",
  "messages": [
    {
      "timestamp": "1741512267722079",
      "device": {
        "id": "852781a2-6167-48e8-a00e-b9c94db50afb",
        "name": "r2d2",
        "protocol": "http"
      },
      "payload": "dGVzdA==",
      "metadata": {
        "Accept": "*/*",
        "Content-Length": "4",
        "Content-Type": "application/x-www-form-urlencoded",
        "Host": "receivers.iot-exchange.io",
        "Location": "/uplink",
        "Method": "POST",
        "User-Agent": "curl/7.68.0"
      }
    }
  ]
}
```
- The payload field is your message encoded as base64 (this allows for binary payloads from your devices)
- For testing you can easily decode base64 on https://www.base64decode.org/


## Need Help?
If you run into any issues, please contact our [support team](mailto:support@iot-exchange.io)

Happy IoT-ing! 🚀

