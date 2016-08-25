# Rock Papper Scissors for Slack

Challenge your Slack team members to a game of Rock Paper Scissors

### TODO
  - Track players wins/losses/ties
  - Scope games/players by team in redis
  - Clean up `challengeHandler` and `interactionHandler`, they're a mess
  - Clean up `utils.go`, it's also a bit sloppy
  - Add bot user to app to send messages `as_user`
  - Add rematch button after game ends
  - Tests! Write tests! Many tests!


### Alternate Workflow Idea
  - `/rps`
    - Story
      - Send button message to channel (visiable to all)
      - Challenger user clicks button, receives ephermeral reply stating their choice
      - Challenger user clicks button(s) again, receives ephemeral reply stating he already has chosen
      - Another user clicks button, button message updates to "@user beat @user, X vs X"
    - Notes
      - Players will consist of the challenger, and the first other user to click button
      - If another user that is not the challenger or first user clicks they will receive an ephemeral message saying too late (or something)

  - `/rps @user`
    - Story
      - Send button message to channel (visiable to all)
      - Same as above, but locked to Challenger and the named user
    - Notes
      - Players will consist of the challenger, and the named user (challengee)
      - If another user clicks a button they will receive an ephemeral message stating that they are not involved in this battle