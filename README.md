# BiliAutoBlacklist

BiliAutoBlacklist is used for automatically blocking users on Bilibili whose usernames contain specific keywords.

You can now use gpt and Feishu open platform to detect and block spam accounts.

## Usage 

Download binary from releases and run.

Using `screen` or something to keep it running.

### Modes 

- `basic`: Keyword Comparison
- `gpt`: Detect using GPT, notify using Feishu when spam accounts are detected.
- `gpt-only`: Detect using GPT, block without confirmation when spam accounts are detected.

## Configuration

When first run, the program will generate a config_example.yaml file.

If using gpt mode, a prompt.example file will also be generated.

Rename `config_example.yaml`,`prompt.example` and to `config.yaml`,`prompt` and edit.

In gpt mode, forward port 9999 to a server with a public ipv4 addr. Start the program.

Then go to Feishu -> Events & Callbacks -> Callback Configuration -> Subscription mode -> Request URL to set the callback url.

```
Example:
127.0.0.1:9999 -> your_server_ip:80
callback url: http://your_server_ip:80/
```
### Configuration Example
```yaml
cookie: "your_cookie_here" # go to browser -> f12 -> network, open a request to "api.bilibili.com"
                           # copy all cookies here
targetUID: "your_uid_here" # your uid, space.bilibili.com/<uid>
timeDelay: 10             # delay between each request, suggest more than 10
cron: "0 0 0/12 * * *"    # optional, no less than 10 minutes
mode: "gpt"               # required, "basic" (Keyword Comparison) or "gpt" () or "gpt-only" ()
BlackListWord:            # user with these words in username will be blocked
  - "word1"               # TODO: users whose username contain these words will be banned directly in gpt mode
  - "word2"
  - "word3"
feishu:   # go to https://open.feishu.cn and create an app & Add bot features to your app
  appId: "your_app_id_here"  # Credentials & Basic Info Page
  appSecret: "your_app_secret_here"
  email: "your_email_here" # your email of feishu, not enterprise email
  token: "your_token_here" # Events and callbacks -> Encryption strategy -> Verification Token
openaiConfig:
  apiBase: "URL_ADDRESS"  # optional, default is https://api.openai.com/v1/
  apiKey: "your_api_key_here" # required
fansCheckPerDay: 10
```



PS: For those who have been harassed by RouZhuangHao account, please contact [zsq001@zsq001.cn](mailto:zsq001@zsq001.cn) for official prompt.
