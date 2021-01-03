import json

def main():
    cmdlist = []
    with open('resources/redis-commands.json', 'r') as f:
        cmds = json.load(f)
        for c, m in cmds.items():
            tmp = dict()
            tmp['name'] = c
            if 'arguments' in m.keys():
                tmp['arguments'] = convert_str_to_list(m['arguments'], 'name')
                tmp['arguments'] = convert_str_to_list(m['arguments'], 'type')

            cmdlist.append(tmp)
    
    new_dict = dict()
    new_dict['redisCommands'] = cmdlist
    with open('resources/redis-commands-golang.json', 'w') as f:
        json.dump(new_dict, f, sort_keys=False)


def convert_str_to_list(arglist, key):
    for arg in arglist:
        if key in arg.keys():
            x = arg.get(key)
            if isinstance(x, str):
                arg[key] = [x]
    return arglist

if __name__ == '__main__':
    main()