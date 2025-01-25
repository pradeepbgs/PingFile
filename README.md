# PingFile CLI

PingFile is a command-line tool that allows you to execute API requests from configuration files defined in JSON, YAML, or PKFILE formats. It helps automate and manage API testing and execution, making it easier to work with various API configurations from a single command.

---

## Features

- Execute API requests from configuration files in multiple formats.
- Supports JSON, YAML, and PKFILE formats.
- Colorful output with status codes, headers, and response body.
- Install the PingFile binary to your system's PATH for easy access.

---

## Installation

### Linux

1. Clone the repository:
   ```bash
   git clone https://github.com/pradeepbgs/pingfile.git
   cd pingfile
   ```
2. Build the binary:
    ```bash 
    go build -o pingfile .
    ```
3. Install the binary globally to your system's PATH:
    ```bash 
    sudo ./pingfile install   
    ```
    This command will move the binary to /usr/local/bin and ensure it's accessible from anywhere in your terminal. Make sure the /usr/local/bin directory is in your PATH.

4. After installation, you can run PingFile from anywhere using the following command:
    ```bash
    pingfile run <file>
    ```
Replace <file> with the path to your configuration file.    

### Windows 
Download/clone the repository and open your CMD

1. navigate to Pingfile Directory
    ```bash
    C:\Users\dell\Desktop\PingFile>
    ```
2. Build the binary:
    ```bash 
    go build -o pingfile.exe
    ```
3. Add Binary to usr/bin
    ```bash
    pingfile install
    OR 
    ./pingfile.exe install
    OR 
    pingfile.exe install
    ```
### 4. Add the binary globally to your system's PATH:
   **Open Environment Variables:**
   * Press `Win + S`, type Environment Variables, and select Edit the system environment variables.
   * In the System Properties window, click the Environment Variables button.

   **Edit the PATH Variable:**
   * Under User Variables (for your user account), locate the `Path` variable and click Edit.
   * Click New and add `C:\Users\dell\bin` to the list.
   * Click OK to save changes and close all dialog boxes.

   **Add the .exe Extension**

  Since pingfile is likely an executable binary, rename it to include the .exe extension:
  ```bash 
  ren C:\Users\dell\bin\pingfile pingfile.exe
  ```
  OR manually go to Users\dell\bin and rename the file pingfile.exe

  **then try this command**

  You should see welcome message
  ```bash
  C:\Users\dell\Desktop\PingFile>pingfile
  Welcome to PingFile!
  Use 'pkfile run <file>' to execute API requests from a file.

  C:\Users\dell\Desktop\PingFile>
  ```



### Usage
After installation, you can run PingFile commands directly from the terminal. Here are the available commands:

`run [file]`

Execute API requests from a configuration file.

**Example**

suppose you made an api endpoint , now you want to test your api if its working or not

getAPI.json
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
 * here you can pass --multithread / -m command 
 * this way you can make multiple req multithreaded

```bash
    pingfile run getAPI.json 

    OR 

    pingfile run getAPI.json PostAPI.json --multithread
```
### Example output
**For a successful request:**

in your terminal
```javascript
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
**A Post Request**

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

### Example output
**For a successful request:**

in your terminal
```javascript
pradeep@pradeep:~/Desktop/PingFile$ pingfile run postapi.yaml 

--------------- >>>>
Running PingFile for: postapi.yaml
<<<<---------------

Status Code: 200 OK

Headers:
  Cache-Control: [no-cache]
  Content-Type: [application/json; charset=utf-8]
  X-Powered-By: [DieselJS]
  Date: [Sat, 25 Jan 2025 08:11:08 GMT]
  Content-Length: [70]

Body:
{"status":true,"data":{"hobby":"coding","name":"pradeep"}}

API request executed successfully!
```

### Here are those tree file formats
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

