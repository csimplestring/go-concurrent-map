package hash

const (
	m              = 249997
	bkdr_seed uint = 131
)

func BKDRHash(str string) uint {
	var hash uint = 0
	for _, s := range []byte(str) {
		hash = hash*bkdr_seed + (uint)(s)
	}

	return hash % m
}
