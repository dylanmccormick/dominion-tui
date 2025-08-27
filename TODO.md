# Server
- Implement connection with TCP and handle sockets?
- I don't think it needs to be UDP because actions in Dominion do not happen that quickly. We can wait for other players. 
    - UDP is good for other types of streaming servers becuase we are OK losing data since we're just resending everything next time (well kind of. I think there is some complex XOR logic to rectify missing data from frames

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
