# membots-blockchain-be

## Build service

    make build

## Run service

    make run

## Stop service

    make stop

## Swagger init

    make swagger

## Create and deploy subgraphs

    Step 1: Access https://thegraph.com/studio/
    Step 2: Connect your metamask wallet
    Step 3: Create and deploy memeception-subgraph
        - Click button "Create a Subgraph"
        - Subgraph name should be "memeception" then click button "Create Subgraph"
        - Copy "AUTHENTICATE IN CLI" command at "AUTH & DEPLOY" section
        - cd subgraph/memeception
        - yarn install
        - Paste and execute "AUTHENTICATE IN CLI" command above
        - graph deploy --studio memeception
    Step 4: Create and deploy swap-v3-subgraph
        - Do the same as the step 3
