import requests
import re
import time
import random
import hashlib

UA = "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
proxies = {
    "http": "http://127.0.0.1:7890",
    "https": "http://127.0.0.1:7890"
}
data = {
    "url": "https://www.youtube.com/watch?v=U9e3MFZI3zE",
    "q_auto": "0",
    "ajax": "1"
}
response = requests.post("https://y2mate.tools/mates/en/analyze/ajax",data=data, proxies=proxies)
print(response.json())
