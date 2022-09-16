import requests
import re
import time
import random
import hashlib

UA = "5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
word = "apple"

resp = requests.get("https://fanyi.youdao.com/", headers={"User-Agent": UA})
OUTFOX_SEARCH_USER_ID = re.match(r'OUTFOX_SEARCH_USER_ID=-?\d+@\d+\.\d+\.\d+\.\d+', resp.headers["Set-Cookie"]).group()

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
print(length, bv, salt, sign, ts)

OUTFOX_SEARCH_USER_ID_NCOO = str(2147483647 * random.random())
___rl__test__cookies = str(int(1000 * time.time()))

headers = {
    'Accept':'application/json',
    "User-Agent": UA,
    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    "Origin": "https://fanyi.youdao.com",
    "Referer": "https://fanyi.youdao.com/",
    "Host": "fanyi.youdao.com",
    "Cookie": "{}; OUTFOX_SEARCH_USER_ID_NCOO={}; ___rl__test__cookies={}".format(OUTFOX_SEARCH_USER_ID, OUTFOX_SEARCH_USER_ID_NCOO, ___rl__test__cookies),
    'Accept': 'application/json, text/javascript, */*; q=0.01',
    "Content-Length": str(length),
    'X-Requested-With':'XMLHttpRequest',
    'Connection':'keep-alive',
}

rep = requests.post("https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule", data=data, headers=headers).json()
print(rep.text)