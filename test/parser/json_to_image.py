import json

import matplotlib.pyplot as plt
import networkx as nx
from networkx.drawing.nx_pydot import graphviz_layout
from networkx.readwrite import json_graph



def generate_graph_from_json(json_data):
    G = nx.DiGraph()
    node_id = 0  # Compteur pour les identifiants uniques des nœuds

    def add_nodes(parent_id, children, depth):
        nonlocal node_id
        for child in children:
            node_type = child['Type']
            current_id = node_id
            G.add_node(current_id, type=node_type, depth=depth,
                       color='skyblue')  # Ajout de l'attribut 'type', 'depth' et color
            if node_type == ',':  # Si le nœud est une virgule, on lui attribue le type 'COMMA'
                G.nodes[current_id]['type'] = ' ,'
            elif node_type == ':':
                G.nodes[current_id]['type'] = '":"'
            G.add_edge(parent_id, current_id)
            node_id += 1
            if child['Children']:
                add_nodes(current_id, child['Children'], depth=depth - 1)
            else:
                G.nodes[current_id]['color'] = 'red'

    root_id = node_id

    G.add_node(root_id, type=json_data['Type'], color='green', depth=0)  # Nœud racine avec son identifiant unique
    node_id += 1
    add_nodes(root_id, json_data['Children'], -1)

    def drawGraph(G, savePath, drawpath):
        plt.figure()

        # Création d'un dictionnaire d'identifiants uniques et de leurs types pour référence
        node_ids = {node: str(node) for node in G.nodes}
        # pos = nx.multipartite_layout(G, subset_key='depth', align='horizontal', scale=40)
        nx.drawing.nx_pydot.write_dot(G, "graph.dot")
        pos = graphviz_layout(G, prog="dot")
        colors = [G.nodes[node]['color'] for node in G.nodes]

        offset = -10
        pos_labels = {}
        keys = pos.keys()
        for key in keys:
            x, y = pos[key]
            pos_labels[key] = (x, y + offset * ((key % 4) - 1.5))
        nx.draw_networkx_nodes(G, pos, node_size=50, node_color=colors)
        nx.draw_networkx_edges(G, pos, arrows=True)
        nx.draw_networkx_labels(G, labels=node_ids, pos=pos_labels, font_color='black', font_size=8, font_weight='bold')

        if savePath:
            with open(savePath, 'w') as f:
                f.write(str(json_graph.node_link_data(G)))

        plt.savefig(drawpath)
        plt.show()


    drawGraph(G,"","parsetree.png")

    def count_terminal_childs(node):
        successors = list(G.successors(node))
        count_terminal = 0
        count_not_terminal = 0
        for successor in successors:
            if G.nodes[successor]['color'] == 'red':
                count_terminal += 1
            else:
                count_not_terminal += 1
        return (count_terminal, count_not_terminal)

    def up_terminal_nodes():
        L = []
        for node in G.nodes:
            c = count_terminal_childs(node)
            if c[0] == 1 and c[1] <= 1 and G.nodes[node]['type'] == '':
                successors = list(G.successors(node))
                for successor in successors:
                    if G.nodes[successor]['color'] == 'red':
                        G.nodes[successor]['color'] = 'skyblue'
                        G.nodes[node]['color'] = 'red'
                        G.nodes[node]['type'] = G.nodes[successor]['type']
                        G.nodes[successor]['type'] = ''

    def delete_nodes():
        L = []
        for node in G.nodes:
            if G.nodes[node]['type'] == '' and len(list(G.successors(node))) == 0:
                L.append(node)
        G.remove_nodes_from(L)

    def delete_chain_nodes():
        continue_delete = True
        i = 0
        while continue_delete:
            print(i)
            continue_delete = False
            L = []
            for node in G.nodes:
                if G.nodes[node]['type'] == '' and len(list(G.successors(node))) == 1 and len(
                        list(G.predecessors(node))) == 1 and G.nodes[list(G.successors(node))[0]]['type'] == '' and \
                        G.nodes[list(G.successors(node))[0]]['color'] == 'skyblue':
                    L.append(node)
                    G.add_edge(list(G.predecessors(node))[0], list(G.successors(node))[0])
                    continue_delete = True
            G.remove_nodes_from(L)
            i+=1

    for j in range(10):
        up_terminal_nodes()
        delete_nodes()
        delete_chain_nodes()



    drawGraph(G,"new_graph.json","ast.png")



# Chargement des données depuis le fichier JSON
# with open('return.json', 'r') as file:
#     json_data = json.load(file)
#
# generate_graph_from_json(json_data)


def get_Types(types,key):
    r = types.get(str(key), None)
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

    print(graph,types)

    def drawGraph(G, savePath, drawpath, id=False):
        plt.figure(figsize=(10, 8))  # Adjust figure size for better visibility

        # Create a dictionary of unique identifiers and types for reference
        if id:
            node_ids = {node: str(node) for node in G.nodes}
        else:
            node_ids = {node: G.nodes[node]['type'] for node in G.nodes}

        # Use a circular layout for better organization
        pos = graphviz_layout(G, prog="dot")
        colors = [G.nodes[node]['color'] for node in G.nodes]

        offset = 40  # Increase the offset for better label positioning
        pos_labels = {key: (x, y + offset) for key, (x, y) in pos.items()}

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
        print(ind,depth)
        if len(graph[str(ind)].keys())==0:
            G.add_node(ind, type=get_Types(types,ind), depth=depth,color='red')
        else:
            G.add_node(ind, type=get_Types(types,ind), depth=depth,color='skyblue')
        for child in graph[str(ind)].keys():
            G.add_edge(ind,child)
            stack.append((child,depth+1))
    print([G.nodes[node]['type'] for node in G.nodes])
    drawGraph(G,"","./test/parser/ast.png")

with open('./test/parser/ast.json', 'r') as file:
    json_data = json.load(file)

gen_graph_jsongraph(json_data)