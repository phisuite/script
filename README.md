# Phi Suite Script

| **Homepage** | [https://phisuite.com][0]        |
| ------------ | -------------------------------- | 
| **GitHub**   | [https://github.com/phisuite][1] |

## Overview

This project contains the **Phi Script Interpreter**.  
**Phi Script** is a markup language to read and write **Phi Suite Schemas**.  
The interpreter connects to a **Phi Suite Kernel** via its **Schema Inspector & Editor**.

## Inspect Existing Schemas

```bash
make inspect INSPECTOR=<inspector address>
```

## Create & Update Schemas

```bash
make update INSPECTOR=<inspector address> EDITOR=<editor address> [file]
```

## Phi Script Example

```
event? RegisterUserRequested:0.1
  string name
  number? age

event! RegisterUserSucceed:0.1
  string id
  string name

event? RegisterUserFailed:0.1
  string message

entity! User:0.1
  string id
  string name
  number? age

entity? Error:0.1
  string message

process~ RegisterUser:0.1
  input RegisterUserRequested:0.1
  output RegisterUserSucceed:0.1
  error RegisterUserFailed:0.1

  input User:0.1
  output User:0.1
  error Error:0.1
```

[0]: https://phisuite.com
[1]: https://github.com/phisuite
[2]: https://github.com/phisuite/schema
