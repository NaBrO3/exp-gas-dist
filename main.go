package main

// generate test data
func main() {
	obj := &exp{}
	obj.reset()
	obj.init(0)
	obj.gen(nov)
	defer obj.close()
}

// number of voters (in worst case 100m)
const nov = 100000000
