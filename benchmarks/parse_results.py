import numpy as numpy


t = open("bench.log","r").readlines()[5:65]

ls = [[tx.strip() for tx in ti.split("\t")] for ti in t]

