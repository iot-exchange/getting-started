# IoT Exchange - Getting Started

Welcome to the IoT Exchange Getting Started repository! 
This guide will help you quickly getting started with testing your device.

## Prerequisites

- An IoT Exchange account ([Sign up here](https://portal.iot-exchange.io))
- `git` + `curl` + `bash` installed to use the scripts provided in this repo

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
- Save the file as `<your device>.pem`.

### 5. Add an Integration
Integrations are external applications that will receive the messages from
all devices in the project.

For testing we are going to add a test integration with the free service
[webhook.site](https://webhook.site) that allows you to receive webhooks in  your browser in realtime.


- Go to https://webhook.site and click on "Your unique URL" to copy the unique url it has generated for you.
- Go back to the IoT Exhchange and navigate to the integrations of your project.
- Click **Add Integration**.
- Paste your unique webhook url into the URI field. You can leave the rest as-is.
- Finish adding the Integration by clicking on **Save**

### 6. Send Data Using the provided Script
This repository contains a script to send data from your device to IoT Exchange.
It uses `curl` to send uplink messages to the iot-exchange as if it was
coming from an actual IoT device.

**Note**: This uses the `https` protocol (with client certificates); for actual IoT devices you want probably want
to use `COAP` or `UDP` with dTLS (with client certificates) as that uses less data and is more suited for IoT networks.

#### Steps to run the script:
1. Clone this repository:
   ```sh
   git clone https://github.com/iot-exchange/getting-started.git
   cd getting-started
   ```
2. Run the script with :
   ```sh
   echo -n "<Your message goes here>" | ./uplink.sh /path/to/<your device>.pem
   ```

### 7. View your message
- Go to https://webhook.site this should now display your message
- It should look something like this:
```json
{
  "timestamp": "1741512267750437",
  "messages": [
    {
      "timestamp": "1741512267722079",
      "device": {
        "id": "<YOUR DEVICE UUID>",
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
- For testing purposes you can easily decode base64 on https://www.base64decode.org/

### 8. Other protocols

Currently we support the following protocols over the internet: 

| Protocol  | Address                   | Port |
|-----------|---------------------------|------|
| https     | receivers.iot-exchange.io | 443  |
| coap+dtls | receivers.iot-exchange.io | 5684 |
| udp+dtls  | receivers.iot-exchange.io | 4433 |

For https and coap+dtls requests to all paths are accepted and directly forwarded to the integration you 
specified for your projects (as are the headers). For udp+dtls you cannot specify the path/headers and the payload
is directly forwarded to your integrations.

For example uses see [our golang client example](./client-golang/)

To run the golang example you'll need to install golang (>=1.23) and then you can run it like this:  
```sh
cd client-golang
echo -n "<Your message goes here>" | go run main.go /path/to/<your device>.pem
```
This will send your message via all 3 protocols.
If all goes well you'll see something like this:

```
2025/04/22 09:05:52 --------------------------------------------------------------------------------------
2025/04/22 09:05:52 => sending message via coap/dtls to receivers.iot-exchange.io:5684
2025/04/22 09:05:52 <= response: Code: Empty, Token: 46d7880f509acc18, ContentFormat: text/plain; charset=utf-8, Type: Acknowledgement, MessageID: 44904
2025/04/22 09:05:52 <= done
2025/04/22 09:05:52 --------------------------------------------------------------------------------------
2025/04/22 09:05:52 => sending message via udp/dtls to receivers.iot-exchange.io:4433
2025/04/22 09:05:52 <= done
2025/04/22 09:05:52 --------------------------------------------------------------------------------------
2025/04/22 09:05:52 => Sending message via http to https://receivers.iot-exchange.io
2025/04/22 09:05:53 <= response: &{Status:201 Created StatusCode:201 Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[Content-Length:[0] Content-Type:[text/plain] Date:[Tue, 22 Apr 2025 09:05:53 GMT] Server:[fasthttp]] Body:{} ContentLength:0 TransferEncoding:[] Close:false Uncompressed:false Trailer:map[] Request:0xc000154640 TLS:0xc000248180}
2025/04/22 09:05:53 <= done
2025/04/22 09:05:53 --------------------------------------------------------------------------------------
```

If you run it from a different directory you might also see a message like this:
```
2025/04/22 09:05:52 [WARNING] unable to load ca.pem; proceeding without checking server certificate
```

This means it could not load the `ca.pem` that it uses to verify the authenticity of the server.
For testing purposes this is harmless; for production code always make sure you verify the server certificates with 
a valid ca.pem 

## Faq

### I get no errors, but I don't see any messages arriving in webhook.site

The free version of https://webhook.site will automatically remove endpoints when they are not used. If you check back 
on https://webhook.site you might get a new endpoint. 
If it has been a while since you created the webhook.site endpoint, please check if it has changed and update the 
endpoint in your project accordingly. 

### In the webhook data the protocol is always the same, regardless of the protocol I use

This is correct. Currently is shows you the metadata of the device and the protocol is the protocol you configured in 
the iot-exchange when adding/editing the device.

This might change in the future.

## Need Help?
If you run into any issues, please contact our [support team](mailto:support@iot-exchange.io)

Happy IoT-ing! 🚀

