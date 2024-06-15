# BiliAutoBlacklist

BiliAutoBlacklist is used for automatically blocking users on Bilibili whose usernames contain specific keywords.

## Usage 

Download binary from releases and run.

Due to the limitation of Bilibili API, the program can only detect latest 50 * 5 fans.

Make sure to keep it running in the background. Using `screen` or something.

## Configuration

When first run, the program will generate a config_example.yaml file.

Rename config_example.yaml to config.yaml and edit.

```yaml
cookie: "your_cookie_here" # go to browser -> f12 -> network, open a request to "api.bilibili.com"
                           # copy all cookie to here
targetUID: "your_uid_here" # your uid, space.bilibili.com/<uid>
timeDelay: 10             # delay between each request, suggest more than 10
cron: "0 0 0/12 * * *"    # optional, no less than 10 minutes
BlackListWord:            # user with these words in username will be blocked
  - "word1"
  - "word2"
  - "word3"
```

