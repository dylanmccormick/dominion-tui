# Server
- Implement connection with TCP and handle sockets?
- I don't think it needs to be UDP because actions in Dominion do not happen that quickly. We can wait for other players. 
    - UDP is good for other types of streaming servers becuase we are OK losing data since we're just resending everything next time (well kind of. I think there is some complex XOR logic to rectify missing data from frames

- 9/1/25
    - Need to update the server so that it uses channels to communicate. Right now we're just sending a message in a loop, but we need it to only send a message when one is needed. Probably use some sort of switch statement to see if there's anything coming from different channels. Maybe I'll start by making some kind of chat feature so different users can send a message. This would be a similar structure to when users are sending input for the game loop. (But that would come from the TUI)

- 9/2/25
    - Need to figure out how to send messages between client and server. Probably need some kind of json schema or something to categorize messages. Now that I have chat working I'm going to have to send a bunch of data over that same connection. 
    - Data Needed:
        - Game state
        - Chat Messages
        - User commands  (settings, requests, etc)  
        - User actions (actual game play actions)
    - I imagine these would all be different channels that write to the same TCP connection pool? Or would this be one handler that handles any type of message? I think I like using multiple channels to do this

- Figure out a way to assign players to rooms and create new rooms
- Assign players an order in the room and allow them to change it???? Maybe later with the changing
- Figure out how to get users to start the game when they are ready
- Users get to pick cards allowed in their deck and the amount (Maybe I'd have to read the rules of dominion)

# Game
- Game loop
- Build out card definitions and actions
- Eventually add a bot to play known strategies (Big money etc)
- 

# Client
- TUI
- Figure out how to render cards and game state from game data
- 


# Other Considerations
- Will I need some kind of in-memory store for game data and such? Would it be better to write some of that to disk since game can be slower?
- User authentication? Don't want this thing to be spammed by scrapers or bots or whatever. I have hundreds of dollars in arbys coupons I don't want stolen


## Important Tips?
- Reading from a TCP connection is just like reading from a file, but you don't get to choose how much data you get. So every message needs to have a set of chars at the end to signal that is the end of the message. This is the importance of "\r\n" in HTTP protocol. I need to be parsing for those breaks, but if I read a lot more data than I expect, I need to save that data for teh next round of message parsing. 
