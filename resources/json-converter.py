import json

def main():
    cmdlist = []
    with open('resources/redis-commands.json', 'r') as f:
        cmds = json.load(f)
        for c, m in cmds.items():
            m['name'] = c
            if 'arguments' in m.keys():
                m['arguments'] = convert_str_to_list(m['arguments'], 'name')
                m['arguments'] = convert_str_to_list(m['arguments'], 'type')

            cmdlist.append(m)
    
    new_dict = dict()
    new_dict['redisCommands'] = cmdlist
    with open('resources/redis-commands-golang.json', 'w') as f:
        json.dump(new_dict, f, indent=4, sort_keys=True)


def convert_str_to_list(arglist, key):
    for arg in arglist:
        if key in arg.keys():
            x = arg.get(key)
            if isinstance(x, str):
                arg[key] = [x]
    return arglist

if __name__ == '__main__':
    main()