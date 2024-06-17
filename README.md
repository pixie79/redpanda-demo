# Redpanda Demo

Development is using Visual Studio Code with the Devcontainer plugin. Once the code is cloned reopen using devcontainers and then run:

```zsh
setup-devcontainer.sh
```

## Demo Environment Setup

Setup the demo / test environment

```zsh
task demo-start
```

Generate dummy data

```zsh
task build-test-data-generator
task generate-test-data
```

This demo masks customers data first and last name fields. By default it does not mask the customer _Jane Smith_ to prove the masking is working on the other records. As the generate-test-data is random it could by chance not create a record for Jane Smith, if so then please check _test-data/demoEvent.json_ (the data source) and update the Environment variable (UNMASKED*CUSTOMERS) at the top of \_Taskfile.yml* to not mask a person from the file of your choice.

Build and Deploy the transform

```zsh
task deploy-demo
```

Build the data loader

```zsh
task build-test-data-loader
```

In a second terminal run one of the following:

```zsh
rpk transform logs demo -f
```

Then in the first terminal run one of the following:

```zsh
task load-td-demoEvent
```

At this point you should check the Redpanda console and look in the _demo_ and _output-demo_ topics. You should see that the source topic _demo_ does not have its data masked whereas all records (except the Jane Smith records) are masked in the _output-demo_ topic.

To clean up the demo

```zsh
task demo-stop
```

## Integration Tests

To debug failing integration tests you can use the rpk tool to interegate the test Redpanda testcontainer instance when running. In order to do this do the following:

```zsh
task create-rpk-test-profile
```

You may need to edit the profile and ports for each test run as these change dynamically. Do this by checking the test output for the IP/Port for each of the services then edit the test profile using the command:

```zsh
rpk profile edit
```

Once the profile is saved the rpk commands will work against the local testcontainer cluster.

## Switching RPK profiles

The rpk tool uses profiles to identify which cluster to operate against, in order to switch profiles run:

```zsh
rpk profile list
rpk profile use *PROFILE_NAME*
```
