import json
from collections import OrderedDict

import matplotlib.pyplot as plt
import networkx as nx
from networkx.drawing.nx_pydot import graphviz_layout
from networkx.readwrite import json_graph

import warnings
warnings.filterwarnings("ignore")

def get_Types(types,key):
    r = types.get(key, "")
    if ":" in r:
        return f'"{r}"'
    elif r == ",":
        return '","'
    else:
        return r

def gen_graph_jsongraph(graphStruct):
    G = nx.DiGraph()
    graph = graphStruct["gmap"]
    types = graphStruct["types"]
    terminals = graphStruct["terminals"]
    meaningful = graphStruct["meaningful"]
    def drawGraph(G, savePath, drawpath, id=False):
        plt.figure(figsize=(len(G.nodes)**0.5,len(G.nodes)**0.5))  # Adjust figure size for better visibility
        # Create a dictionary of unique identifiers and types for reference
        if id:
            node_ids = {node: f"{node} : {G.nodes[node]['type']}" for node in G.nodes}
        else:
            node_ids = {node: G.nodes[node]['type'] for node in G.nodes}

        # Use a circular layout for better organization
        pos = graphviz_layout(G, prog="dot")
        colors = [G.nodes[node]['color'] for node in G.nodes]

        offset = 3  # Increase the offset for better label positioning
        pos_labels = {key: (x, y + offset * (((x+y)%7)-3)) for key, (x, y) in pos.items()}

        # Adjust node size and edge width for better visibility
        nx.draw_networkx_nodes(G, pos, node_size=100, node_color=colors, alpha=0.8)
        nx.draw_networkx_edges(G, pos, arrows=True, width=1.0, alpha=0.5)

        # Adjust font size and weight for better label readability
        nx.draw_networkx_labels(G, labels=node_ids, pos=pos_labels, font_color='black', font_size=6, font_weight='bold')

        if savePath:
            with open(savePath, 'w') as f:
                f.write(str(json_graph.node_link_data(G)))

        plt.savefig(drawpath, format="png", dpi=300)

    stack = [(0,0)]

    while stack:
        ind,depth = stack.pop()
        if ind in terminals.keys():
            G.add_node(ind, type=get_Types(types,ind), depth=depth,color='red')
        elif ind in meaningful.keys():
            G.add_node(ind, type=get_Types(types,ind), depth=depth,color='orange')
        else:
            G.add_node(ind, type=get_Types(types,ind), depth=depth,color='skyblue')
        for child in graph[ind]:
            G.add_edge(ind,child)
            stack.append((child,depth+1))
    drawGraph(G,"","./test/parser/ast.png",True)

def parse_int_keys(pairs):
    result = OrderedDict()
    for key, value in pairs:
        if key == "gmap":
            result[key] = OrderedDict(sorted((int(k), sorted(v)) for k, v in value.items()))
        elif key[0] in "0123456789":
            result[int(key)] = value
        else:
            result[key] = value
    return result

with open('./test/parser/ast.json', 'r') as file:
    json_data = json.load(file, object_pairs_hook=parse_int_keys)

gen_graph_jsongraph(json_data)