import re, os

reply_re = "@[a-z-]+reply"


def return_type(docu: str) -> str:
    match = re.search(reply_re, docu)
    # print(match)
    if match:
        ret = match.group(0)
        return ret
    return None


def main():
    docs = os.listdir('redis-doc/commands')
    for d in docs:
        with open('redis-doc/commands/' + d) as f:
            docu = f.read()
            ret = return_type(docu)
            print("File: {0:25} \t\t Return Type: {1}".format(d, ret))

if __name__ == '__main__':
    main()