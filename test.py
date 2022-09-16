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
print(length, bv, salt, sign, ts)

OUTFOX_SEARCH_USER_ID_NCOO = str(2147483647 * random.random())
___rl__test__cookies = str(int(1000 * time.time()))

headers3 = {
    "User-Agent": UA2,
    # "Origin": "https://fanyi.youdao.com",
    "Referer": "https://fanyi.youdao.com/",
    "Host": "rlogs.youdao.com",
    "Cookie": "EARCH_USER_ID={}; OUTFOX_SEARCH_USER_ID_NCOO={}".format(OUTFOX_SEARCH_USER_ID, OUTFOX_SEARCH_USER_ID_NCOO),
}

headers4 = {
    "User-Agent": UA2,
    "Origin": "https://fanyi.youdao.com",
    "Referer": "https://fanyi.youdao.com/",
    "Host": "fanyi.youdao.com",
    "Cookie": "EARCH_USER_ID={}; OUTFOX_SEARCH_USER_ID_NCOO={}; ___rl__test__cookies={}".format(OUTFOX_SEARCH_USER_ID, OUTFOX_SEARCH_USER_ID_NCOO, ___rl__test__cookies),
}

# 注册cookie?
response1 = requests.get("https://rlogs.youdao.com/rlog.php?_npid=fanyiweb&_ncat=event&_ncoo={}&_nssn=NULL&_nver=1.2.0&_ntms={}&_nhrf=newweb_translate_text".format(OUTFOX_SEARCH_USER_ID_NCOO, ___rl__test__cookies), headers=headers3)
print("response1", response1)

headers = {
    "User-Agent": UA,
    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    "Origin": "https://fanyi.youdao.com",
    "Referer": "https://fanyi.youdao.com/",
    "Host": "fanyi.youdao.com",
    "Cookie": "EARCH_USER_ID={}; OUTFOX_SEARCH_USER_ID_NCOO={}; ___rl__test__cookies={}".format(OUTFOX_SEARCH_USER_ID, OUTFOX_SEARCH_USER_ID_NCOO, ___rl__test__cookies),
    'Accept': 'application/json, text/javascript, */*; q=0.01',
    "Content-Length": str(length)
}

response2 = requests.get("https://fanyi.youdao.com/ctlog?pos=undefined&action=&sentence_number=1&type=en2zh-CHS", headers=headers4)
print("response2", response2)
response3 = requests.get("https://fanyi.youdao.com/ctlog?pos=undefined&action=RESULT_DICT_SHOW", headers=headers4)
print("response3", response3)

headers2 = {
    'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
    'Cookie': 'OUTFOX_SEARCH_USER_ID=-941548234@10.108.162.139; OUTFOX_SEARCH_USER_ID_NCOO=194331207.38996145; ___rl__test__cookies=1663287302298',
    'Origin': 'https://fanyi.youdao.com',
    'Referer': 'https://fanyi.youdao.com/',
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36',
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

rep = requests.post("https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule", data=data, headers=headers)
print(rep.text)