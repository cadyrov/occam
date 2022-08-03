# Occam.fi test task

## env settings
### default params palced in .config.json
### all available flags passes throught -- flags
### ENV can be set with OCCAM prefix with structure like config json
    ``` example log level : OCCAM_PROJECT_LOG_LEVEL=debug```

### algoritm realized with moving average 

## for test origins mocked in service in providers/origin/mock

###  import paremeters 
```
OCCAM_PROJECT_SHIFT - define parameter  shift of moving average
OCCAM_PROJECT_OUTPUT - define main output. by default - std_out
OCCAM_PROJECT_PRECISION_SECOND - time ticker to get result. by default equal 60
```

###  build 
```
    make build
```

###  run
```
    make run
```