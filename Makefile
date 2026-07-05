trap:
	go run . maps/trap.map a z 5
2:
	go run . maps/london.map waterloo st_pancras 2
3:
	go run . maps/london.map waterloo st_pancras 3
4:
	go run . maps/london.map waterloo st_pancras 4
5:
	go run . maps/london.map waterloo st_pancras 100
6:
	go run . maps/london.map waterloo st_pancras 1
7:
	go run . maps/london.map waterloo st_pancras 4
8:
	go run . maps/8.map bond_square space_port 4 | wc -l
9:
	go run . maps/9.map jungle desert 10 | wc -l
10:
	go run . maps/10.map beginning terminus 20 | wc -l
11:
	go run . maps/11.map two four 4 | wc -l
12:
	go run . maps/12.map beethoven part 9| wc -l
13:
	go run . maps/13.map small large 9 | wc -l
14:
	go run . maps/london.map a b 
15:
	go run . maps/london.map a b 1 hello
