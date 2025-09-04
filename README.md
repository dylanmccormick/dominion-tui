# Dominion TUI

Objectives
- Learn how to build TUIs with Bubbletea or Opentui (undecided) 
- Learn how to create a server that can handle connections from clients which I made
- Allow multiple games at once over TCP/UDP (undecided. Probably TCP)

### Idea:
Running the dominion executable takes you to the home page for the TUI. From here, you are able to set a username and decide if you want to join a server or host a server.

#### Hosting:
If you are the host of the server, your IP address will be printed to the console and you will be told to share that with a friend. The dominion config file will figure out which port to open to communicate with the other users. 

<!-- I'll have to figure out network security -->
#### Joining:
Just put in the IP address of your friend's server and then you'll be off to the races.

All commands in the game will be able to be done by keyboard shortcuts or prompted options for number selection with arrow keys, keyboard input (type a number), or vim keybinds (HJKL)


### Next Steps:

Idk some design in figma maybe?

Play with bbt or get the game loop. Many directions this can go


## Message Protocol
We will need a message protocol to help facilitate and organize messages through the system. There will be multiple types of messages and each message will have an ID. 

prompt: This is any message that prompts the client for input. Will be used for game actions and login and maybe other things
```json
{
    "version": 1,
    "message_id": "auth_001",
    "type": "prompt",
    "ack_needed": true,
    "body": {
        "prompt_type": "authentication",
        "title": "Login Required",
        "options": ["login", "register"]
        "fields": [
            {"name": "username", "type": "text", "required": true},
            {"name": "password", "type": "password", "required": true}
        ]
    }
}
```

prompt_response: this is a response to a prompt. We'll need to know which prompt we are responding to
```json
{
    "version": 1,
    "message_id": "auth_002",
    "type": "prompt_response",
    "ack_needed": true,
    "body": {
        "prompt_id": "auth_001",
        "action": "login"
        "data": {
            "username": "Alice",
            "password": <hashed_password>
        }
        "fields": [
            {"name": "username", "type": "text", "required": true},
            {"name": "password", "type": "password", "required": true}
        ]
    }
}
```
