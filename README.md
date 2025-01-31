# PingFile CLI

```ocaml
pradeep@pradeep:~/Desktop/PingFile$ pingfile run getAPI.json

--------------- >>>>
Running PingFile for: getAPI.json
<<<<---------------

Status Code: 200 OK

Headers:
  Cache-Control: [no-cache]
  Content-Type: [text/plain; charset=utf-8]
  X-Powered-By: [DieselJS]
  Date: [Sat, 25 Jan 2025 08:00:10 GMT]
  Content-Length: [12]

Body:
Hello World!

API request executed successfully!
```
PingFile is a command-line tool that allows you to execute API requests from configuration files defined in JSON, YAML, or PKFILE formats. It helps automate and manage API testing and execution, making it easier to work with various API configurations from a single command.

---

## Features

- Execute API requests from configuration files in multiple formats.
- Supports JSON, YAML, and PKFILE formats.
- Colorful output with status codes, headers, and response body.
- Install the PingFile binary to your system's PATH for easy access.

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/pradeepbgs/PingFile.git
   cd pingfile
   ```
2. Go to bin directory:
    ```bash
    pradeep@pradeep:~/Desktop/PingFile$ ls
    api.json    cmd            go.sum    postapi.yaml
    bin         cookie.pkfile  internal  post.pkfile
    build.sh    cookie.yaml    main.go   README.md
    CMakeFiles  go.mod         pingfile  root.pkfile
    pradeep@pradeep:~/Desktop/PingFile$ cd bin/
    pradeep@pradeep:~/Desktop/PingFile/bin$ 
    ```
   
3. Install the binary globally to your system's PATH:

    **Based on your Operating System**

    1. **For Linux**
        ```bash 
        sudo ./pingfile-linux install

        OR ARM Based CPU
        sudo /.pingfile-linux-arm install   
        ```
        This command will move the binary to /usr/local/bin and ensure it's accessible from anywhere in your terminal. Make sure the /usr/local/bin directory is in your PATH.
    2. **For Windows**
        ```bash
         ./pingfile-windows.exe install
         ```
         1. Add the binary globally to your system's PATH:
        
            **Open Environment Variables:**
            * Press `Win + S`, type Environment Variables, and select Edit the system environment variables.
            * In the System Properties window, click the Environment Variables button.

        2. Edit the PATH Variable:
            * Under User Variables (for your user account), locate the `Path` variable and click Edit.
            * Click New and add `C:\Users\dell\bin` to the list.
            * Click OK to save changes and close all dialog boxes.

        3. Add the .exe Extension

        **Since pingfile is likely an executable binary, rename it to include the .exe extension:**
        ```bash 
        ren C:\Users\dell\bin\pingfile pingfile.exe
        ```
        OR manually go to Users\dell\bin and rename the file pingfile.exe
    3. **For Macos**
        ```bash 
        sudo ./pingfile-macos install

        OR ARM Based CPU
        sudo /.pingfile-macos-arm install   
        ```
        **i dont know how to set path in macOS as i dont have but you know the process.**
4. **After installation, you can run PingFile from anywhere using the following command:**
    ```bash
    C:\Users\dell\Desktop\PingFile>pingfile
    Welcome to PingFile!
    Use 'pingfile run <file>' to execute API requests from a file.

    C:\Users\dell\Desktop\PingFile>
    ```



### Usage
After installation, you can run PingFile commands directly from the terminal. Here are the available commands:

`run [file]`

Execute API requests from a configuration file.

**Example**

suppose you made an api endpoint , now you want to test your api if its working or not

**A normal GET Request**

getrequest.json
```json
{
    "name":"Ping hello world /",
    "url":"http://localhost:3000/",
    "headers":{
        "Method":"GET"
    }
}
```
### Run the command
**Note**
 * here you can pass --multithread / -m flag for multithreaded running apis 
 * you can pass --save / -s to save the response

```bash
    pingfile run getrequest.json 

    OR 

    pingfile run getrequest.json postrequest.json --multithread -s
```
### Example output
**For a successful request:**

in your terminal
```javascript
pradeep@pradeep:~/Desktop/PingFile$ pingfile run getrequest.json

--------------- >>>>
Running PingFile for: getrequest.json
<<<<---------------

Status Code: 200 OK

Headers:
  Cache-Control: [no-cache]
  Content-Type: [text/plain; charset=utf-8]
  X-Powered-By: [DieselJS]
  Date: [Sat, 25 Jan 2025 08:00:10 GMT]
  Content-Length: [12]

Body:
Hello World!

API request executed successfully for getrequest.json
```

### 1. Post Request with files
postwithfiles.json
```json
{
    "name": "POST req to /body with files",
    "url": "http://localhost:3000/body",
    "filePath": "pkfile/postwithfile.json",
    "saveResponse": false,
    "includeCookie": true,
    "headers": {
        "Method": "POST"
    },
    "body": {
        "name": "pradeep",
        "password": "password hi hai",
    },
    "file": [
            {
                "name": "avatar",
                "path": "/home/pradeep/Pictures/project1.png"
            },
            {
                "name":"avatar2",
                "path":"/home/pradeep/Pictures/project1.png"
            }
        ]
}
```
### 2. A Noraml Post Request

postapi.yaml

```yaml
name: POST request
url: http://localhost:3000/body
headers:
  Method: POST
  Content-Type: application/json
body:
  name: pradeep
  hobby: "coding"
```

## How you should configurate your file?
**Here is an Example with .json file**
```json
{
    "name":"Ping hello world /",
    "filePath":"pkfile/getApi.json",
    "saveResponse":true,
    "includeCookie":true,
    "includeCredentials":false,
    "url":"http://localhost:3000/",
    "headers":{
        "Method":"GET"
    },
    "credentials": {
        "type": "basic",
        "username": "${API_USERNAME}",
        "password": "${API_PASSWORD}"
    }
}
```
**You can write simillar config for .yaml file aslo**


## Here are those tree file formats
**pkfile**

it's just a json but with pkfile extension , in this extension you will get snippets
```json
{
    "name":"POST req to /body",
    "url":"http://localhost:3000/body",
    "headers":{
        "Method":"POST",
        "Content-Type":"application/json"
    },
    "body":{
        "name":"pradeep",
        "password":"password hi hai"
    }
}
```

**json**
```json
{
    "name":"Ping hello world /",
    "url":"http://localhost:3000/",
    "headers":{
        "Method":"GET"
    }
}
```
**yaml**
```yaml
name: POST request
url: http://localhost:3000/body
headers:
  Method: POST
  Content-Type: application/json
body:
  name: pradeep
  hobby: "coding"
```
