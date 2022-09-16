import requests
import re
import time
import random
import hashlib

UA = "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
UA2 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
word = "apple"

resp = requests.get("https://fanyi.youdao.com/", headers={"User-Agent": UA2})
print(resp.headers)
OUTFOX_SEARCH_USER_ID = re.match(r'OUTFOX_SEARCH_USER_ID=-?\d+@\d+\.\d+\.\d+\.\d+', resp.headers["Set-Cookie"]).group()
print(OUTFOX_SEARCH_USER_ID)

t = hashlib.md5(UA.encode("utf-8")).hexdigest()
r = int(1000 * time.time())
i = str(r) + str(random.randint(0,9))

ts = r
bv = t
salt = i
sign = hashlib.md5(("fanyideskweb" + word + i + "Ygy_4c=r#e#4EX^NUGUc5").encode("utf-8")).hexdigest()

data = {
    "i": word,
    "from": "AUTO",
    "to": "AUTO",
    "smartresult": "dict",
    "client": "fanyideskweb",
    "salt": salt,
    "sign": sign,
    "lts": str(ts),
    "bv": bv,
    "doctype": "json",
    "version": "2.1",
    "keyfrom": "fanyi.web",
    "action": "FY_BY_REALTlME"
}
length = -1
for key, value in data.items():
    length += len(key) + len(value) + 2
print(length, bv, salt, sign)

headers = {
    "User-Agent": UA,
    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    "Origin": "https://fanyi.youdao.com",
    "Referer": "https://fanyi.youdao.com/",
    "Host": "fanyi.youdao.com",
    "Cookie": "EARCH_USER_ID={}; OUTFOX_SEARCH_USER_ID_NCOO={}; ___rl__test__cookies={}".format(OUTFOX_SEARCH_USER_ID, 2147483647 * random.random(), int(1000 * time.time())),
    'Accept': 'application/json, text/javascript, */*; q=0.01',
    'Accept-Encoding': 'gzip, deflate, br',
    'Accept-Language': 'zh-US,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6,ja-CN;q=0.5,ja;q=0.4',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
    'Content-Length': str(length),
    'X-Requested-With': 'XMLHttpRequest'
}

headers2 = {
    'Accept': 'application/json, text/javascript, */*; q=0.01',
    'Accept-Encoding': 'gzip, deflate, br',
    'Accept-Language': 'zh-US,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6,ja-CN;q=0.5,ja;q=0.4',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive',
    'Content-Length': '239',
    'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
    'Cookie': 'OUTFOX_SEARCH_USER_ID=-941548234@10.108.162.139; OUTFOX_SEARCH_USER_ID_NCOO=194331207.38996145; ___rl__test__cookies=1663287302298',
    'Host': 'fanyi.youdao.com',
    'Origin': 'https://fanyi.youdao.com',
    'Pragma': 'no-cache',
    'Referer': 'https://fanyi.youdao.com/',
    'sec-ch-ua': '"Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': "Windows",
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36',
    'X-Requested-With': 'XMLHttpRequest'
}

data2 = {
    "i": word,
    "from": "AUTO",
    "to": "AUTO",
    "smartresult": "dict",
    "client": "fanyideskweb",
    "salt": "16632870367050",
    "sign": "16a61b0f86c272e5ca5ef40f0a59cf68",
    "lts": "1663287036705",
    "bv": "47edca4d7e6ec9bf4fca7156ea36b8ef",
    "doctype": "json",
    "version": "2.1",
    "keyfrom": "fanyi.web",
    "action": "FY_BY_REALTlME"
}

rep = requests.post("https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule", json=data, headers=headers).json()
print(rep)