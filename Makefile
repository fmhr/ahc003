build: main.go
	go build

toolrun:
	go build
	cd tools &&\
		cargo run --release --bin tester in/0000.txt ../solver > out.txt

toolrunbin:
	go build
	./tools/target/release/tester tools/in/0000.txt ./solver > out.txt

localrun:
	go build
	./solver -local < tools/in/0000.txt > out.txt

# example:
# 	./src/solver < tools/example.in > example.out
# 
# vis:
# 	cd tools &&\
# 	cargo run --release --bin vis example.in example.out
# 
# 
# buildcmd:
# 	cd script &&\
# 		go build
# 
# cmdtest:
# 	make buildcmd
# 	./script/script
# 
# 
# pprof:
# 	.solver -cpuprofile cpu.prof < tools/in/0000.txt
# 	pprof -http=localhost:8080 cpu.prof
