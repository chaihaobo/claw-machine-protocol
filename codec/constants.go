package codec

// IndexBoxHost 盒子主机
// IndexMainboard 娃娃机主板
const (
	IndexBoxHost   Index = 0x01
	IndexMainboard Index = 0x02
)

const (
	FrameHead = 0xaa
	FrameTail = 0xdd
)

var (
	ValidIndexes = []Index{IndexBoxHost, IndexMainboard}
)

type Index byte

func IsValidIndex(index Index) bool {
	for _, v := range ValidIndexes {
		if v == index {
			return true
		}
	}
	return false
}
