No one is using this so i don't get any joy at making issue and solving them.
maybe when i want to code in Go again so i wll look into it ( obviously one of the best project of mine )
# PingFile CLI

PingFile is a command-line tool that allows you to execute API requests from configuration files defined in JSON, YAML, or PKFILE formats. It helps automate and manage API testing and execution, making it easier to work with various API configurations from a single command.

---

## Features

- Execute API requests from configuration files in multiple formats.
- Supports JSON, YAML, and PKFILE formats.
- Supports multiple API requests in a single configuration file.
- Colorful output with status codes, headers, and response body.
- Install the PingFile binary to your system's PATH for easy access.

---

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/pradeepbgs/PingFile.git
   cd pingfile
   ```

2. Go to the `bin` directory:

   ```bash
   cd bin
   ```

3. Install the binary globally to your system's PATH:

### For Linux
```bash
chmod +x ./pingfile-linux
sudo ./pingfile-linux install
```

### For Windows
```bash
./pingfile-windows.exe install
```
Add `C:\Users\dell\bin` to your system's PATH for easy access.

### For macOS
```bash
chmod +x ./pingfile-macos
sudo ./pingfile-macos install
```

4. After installation, run PingFile with the following command:
```bash
pingfile
```

---

## Usage

### Single API Request

**getrequest.json**
```json
{
    "name":"Ping hello world /",
    "url":"http://localhost:3000/",
    "headers":{
        "Method":"GET"
    }
}
```

Run the command:
```bash
pingfile run getrequest.json
```

---

### Multiple API Requests in One Configuration File

**multi_apis.json**
```json
{
    "name": "Global api/v1/users route",
    "description": "This is a longer description for the API route.",
    "version": "1.0.0",
    "baseUrl": "http://localhost:3001",
    "apis": [
        {
            "name": "Register",
            "run": true,
            "filePath": "zpkfile/register_response.json",
            "saveResponse": true,
            "includeCookie": false,
            "url": "/api/v1/user/register",
            "headers": { "Method": "POST" },
            "body": {
                "fullname":"pradeep kumar",
                "username":"cityofdeadpeople",
                "email": "kumarpradeepbgs@gmail.com",
                "password": "exvillager"
            }
        },
        {
            "name": "Login",
            "run": true,
            "filePath": "zpkfile/login_response.json",
            "saveResponse": true,
            "includeCookie": true,
            "url": "/api/v1/user/login",
            "headers": { "Method": "POST" },
            "body": {
                "email": "kumarpradeepbgs@gmail.com",
                "password": "linux"
            }
        },
        {
            "name": "Upload video",
            "filePath": "pkfile/getApi.json",
            "run": false,
            "saveResponse": false,
            "includeCookie": true,
            "url": "/api/v1/video/upload",
            "headers": { "Method": "POST" },
            "body": {
                "title": "Proompted Kiddies Learning The Hard Way",
                "description": "Twitch / theprimeagen \nDiscord / discord"
            },
            "file": [
                { "name": "video", "path": "/home/pradeep/Downloads/video.mp4" },
                { "name": "thumbnail", "path": "/home/pradeep/Downloads/image.jpg" }
            ]
        }
    ]
}
```

Run the command:
```bash
pingfile run multi_apis.json
```

---

### Notes
- Use `--multithread` for concurrent requests.
- Use `--save` to save the response to the specified file path.

---

## Example Output
```
pradeep@pradeep:~/Desktop/PingFile$ pingfile run apis.json

--------------- >>>>
Running PingFile for: apis.json
<<<<---------------

Status Code: 200 OK

Headers:
  Cache-Control: [no-cache]
  Content-Type: [text/plain; charset=utf-8]
  X-Powered-By: [DieselJS]

Body:
Hello World!

API request executed successfully!
```

---

## Contribution
Feel free to contribute by opening issues or submitting pull requests.

---

## License
This project is licensed under the MIT License.

