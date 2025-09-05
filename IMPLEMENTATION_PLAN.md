# Dominion TUI Implementation Plan

## File Structure Architecture


```
dominion-tui/
├── cmd/
│   └── dominion-tui/
│       └── main.go                 # Entry point, CLI arg parsing
├── internal/
│   ├── tui/                        # TUI Layer
│   │   ├── app.go                  # Main TUI application controller  
│   │   ├── models/                 # Bubbletea models for each screen
│   │   │   ├── lobby.go           # Lobby/room selection screen
│   │   │   ├── game.go            # Main game screen
│   │   │   └── menu.go            # Main menu screen
│   │   ├── components/             # Reusable UI components
│   │   │   ├── card.go            # Card rendering component
│   │   │   ├── chat.go            # Chat component
│   │   │   └── player_status.go   # Player info display
│   │   └── styles/                # Visual styling
│   │       └── styles.go          # Color schemes, borders, layouts
│   ├── network/                    # Network Layer
│   │   ├── client.go              # TCP client with reliability
│   │   ├── protocol.go            # Message definitions and parsing
│   │   └── connection.go          # Connection management
│   ├── game/                       # Game Layer  
│   │   ├── engine.go              # Core game state management
│   │   ├── player.go              # Player state and actions
│   │   ├── turn.go                # Turn management logic
│   │   ├── rules.go               # Game rules validation
│   │   └── diff.go                # State diffing and patching
│   ├── cards/                      # Data Layer (existing)
│   │   └── cards.go               # Card definitions
│   └── utils/                      # Utilities (existing)
│       └── utils.go               # Helper functions
├── server/                         # Server Package (existing)
│   ├── server.go                  # Enhanced server with new protocol
│   ├── room.go                    # Room management (extract from server.go)
│   ├── user.go                    # User management (existing)
│   └── message.go                 # Enhanced message handling
├── go.mod
├── go.sum
└── README.md
```

## Implementation Plan by Layer

### Phase 1: Foundation Layer (Data + Utils)

**Data Layer Steps:**
1. Enhance card definitions with complete Dominion card set
2. Add card effect descriptions and action types
3. Create card validation functions
4. Add deck building utilities

**Utils Layer Steps:**
1. Add UUID generation functions (replace external dependency)
2. Create JSONPath parsing utilities using string manipulation
3. Add message serialization helpers
4. Create validation helper functions

### Phase 2: Network Layer

**Protocol Design Steps:**
1. Define message types enumeration for all game actions
2. Create message structure with ID, timestamp, type, body
3. Design ACK response format
4. Create message validation functions

**Connection Management Steps:**
1. Build reliable TCP client with connection pooling
2. Implement message queuing and retry logic
3. Add timeout handling and reconnection
4. Create heartbeat/keepalive mechanism

**Message Processing Steps:**
1. Build message parser with delimiter handling
2. Create pending message tracking system
3. Implement acknowledgment waiting with timeouts
4. Add message ordering and duplicate detection

### Phase 3: Game Layer

**State Management Steps:**
1. Design core game state structures (players, supply, turn info)
2. Create state serialization and deserialization
3. Build state validation functions
4. Add state persistence utilities

**Diffing System Steps:**
1. Implement state comparison using reflection
2. Create diff operation structures (set, add, remove)
3. Build JSONPath-style patching system
4. Add diff compression and optimization

**Game Engine Steps:**
1. Create turn-based state machine
2. Implement action validation logic
3. Add card effect resolution system
4. Build victory condition checking

**Rules Engine Steps:**
1. Create action legality checking
2. Implement timing rules (action/buy phases)
3. Add card interaction validation
4. Build game setup and initialization

### Phase 4: Server Layer Enhancement

**Room Management Steps:**
1. Extract room logic from server.go into room.go
2. Add game state persistence per room
3. Implement player turn management
4. Create spectator support

**Message Handling Steps:**
1. Enhance message routing by type
2. Add broadcast messaging to room members
3. Implement state diff distribution
4. Create error handling and recovery

**Game Coordination Steps:**
1. Add turn synchronization across clients
2. Implement action validation before state changes
3. Create game event logging
4. Add reconnection state recovery

### Phase 5: TUI Layer

**Application Structure Steps:**
1. Create main TUI application controller
2. Implement screen routing and navigation
3. Add global state management for UI
4. Create event handling system

**Screen Models Steps:**
1. Build main menu model with server/client options
2. Create lobby model for room selection and chat
3. Implement game screen model with all game elements
4. Add settings and help screen models

**Component System Steps:**
1. Create card rendering component with ASCII art
2. Build chat component with message history
3. Implement player status displays
4. Add supply pile visualization

**Input Handling Steps:**
1. Create keyboard shortcut system
2. Implement vim-style navigation
3. Add number selection for card choices
4. Create command input parsing

### Phase 6: Integration Layer

**Client-Server Integration Steps:**
1. Connect TUI events to network messages
2. Implement game state synchronization
3. Add real-time update handling in UI
4. Create error display and recovery

**Message Flow Steps:**
1. Route user actions through validation
2. Send actions to server with ACK waiting
3. Receive state diffs and apply to local state
4. Update UI components based on new state

**Testing Integration Steps:**
1. Create automated client simulation
2. Add message flow testing
3. Implement load testing with multiple clients
4. Create game logic unit tests

### Phase 7: Polish and Features

**Advanced Features Steps:**
1. Add game replay system using event log
2. Implement AI player support
3. Create game statistics tracking
4. Add custom card set support

**User Experience Steps:**
1. Add smooth animations for card movements
2. Implement sound effects (terminal bell)
3. Create help system with rule explanations
4. Add colorblind-friendly themes

**Reliability Steps:**
1. Implement automatic reconnection
2. Add graceful degradation for network issues
3. Create backup state recovery
4. Add comprehensive error logging

## Implementation Priority Order

1. **Start with Phase 2 (Network)** - Build the reliable messaging foundation
2. **Move to Phase 3 (Game)** - Create the core game logic and state management  
3. **Enhance Phase 4 (Server)** - Upgrade your existing server with new capabilities
4. **Build Phase 5 (TUI)** - Create the user interface layer
5. **Connect Phase 6 (Integration)** - Wire everything together
6. **Polish Phase 7** - Add advanced features and improve UX

Each phase builds on the previous ones, and you can test incrementally as you go. The network layer is crucial to get right first since everything else depends on reliable communication.

## Design Patterns Used

### Layered Architecture
```
┌─────────────────┐
│   TUI Layer     │  ← User interaction, rendering, input handling
├─────────────────┤
│  Network Layer  │  ← TCP communication, message serialization
├─────────────────┤
│   Game Layer    │  ← Business logic, rules, state management
├─────────────────┤
│   Data Layer    │  ← Card definitions, persistence
└─────────────────┘
```

### Key Concepts
- **JSONPath State API**: Internal REST-like API for game state manipulation
- **Event Sourcing**: Game actions as events that can be replayed
- **State Diffing**: Send only changes instead of full state for efficiency
- **Reliable Messaging**: Application-level ACKs for guaranteed delivery
- **Turn-based Synchronization**: Server authority with client validation
