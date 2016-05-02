import csv
import sys

ifile = open(sys.argv[1])
ofile = open(sys.argv[2],"w")

rows = []

users = []

reader = csv.reader(ifile)
for row in reader :
    rows.append(tuple(row))
    if row[0] not in users :
        users.append(row[0])

if len(users) != 2 :
    raise Exception("improper number of users")

umap = {users[0]:users[1],users[1]:users[0]}

writer = csv.writer(ofile)
for row in rows :
    out = (row[0],umap[row[0]],row[1],row[2],row[3])
    writer.writerow(out)



