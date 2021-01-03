import json
import sys
import re
from collections import OrderedDict

commands_json = "./redis-doc/commands.json"
commands_folder = "./redis-doc/commands/"

def main() -> list:
    cmdlist = []
    with open(commands_json) as cj:
        commands = json.load(cj)
        for cmd, dets in commands.items():
            comm = create_command(cmd, dets)
            cmdlist.append(comm)
    
    return cmdlist


def create_command(cmd: str, dets: dict) -> dict:
    documentation = ""
    try:
        cmdfile = cmd.lower().replace(" ", "-") + ".md"
        docfile: str = commands_folder + cmdfile
        with open(docfile) as f:
            documentation = f.read()
    except FileNotFoundError:
        print("File not found: {}".format(docfile))
    
    dets["name"] = cmd
    dets["document"] = documentation
    ret = extract_return(doc=documentation)
    if ret is not None:
        dets["return"] = ret.replace("@", "")
    if 'arguments' in dets.keys():
        args = dets['arguments']
        compat_args = []
        for a in args:
            a2 = compat(a)
            compat_args.append(a2)
        dets['arguments'] = compat_args

    return dets


def compat(a: dict) -> dict:
    namelist = []
    typelist = []
    c = {}
    if 'name' in a.keys():
        n = a['name'] 
        if isinstance(n, str):
            namelist.append(n)
        else:
            namelist = n

    if 'type' in a.keys():
        t = a['type']
        if isinstance(t, str):
            typelist.append(t)
        else:
            typelist = t
    c['name'] = namelist
    c['type'] = typelist

    for k, v in a.items():
        if k != 'name' and k != 'type':
            c[k] = v

    return c
    

def extract_return(doc: str) -> str:
    reply_re = "@[a-z-]+reply"
    match = re.search(reply_re, doc)
    # print(match)
    if match:
        ret = match.group(0)
        return ret
    return None


if not(sys.version_info.major >= 3 and sys.version_info.minor >= 5):
    print("Can run this only for python>=3.5")
    sys.exit(1)

if __name__ == '__main__':
    out = main()
    order = ["name", "summary", "arguments", "return", "complexity", "group", "since", "document"]
    ordered = [OrderedDict(sorted(item.items(), key = lambda item: order.index(item[0]))) for item in out]
    with open('resources/commands-golangcompat.json', 'w') as f:
        json.dump(obj=ordered, fp=f, indent=2, sort_keys=False)

