package utils

import (
	"crypto/sha256"
	"errors"
	"net/url"
)

const (
	digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	pow16 = [10]int{68719476736, 4294967296, 268435456, 16777216, 1048576, 65536, 4096, 256, 16, 1}
	pow62 = [7]int{56800235584, 916132832, 14776336, 238328, 3844, 62, 1}
)

var (
	ErrorHashCollision = errors.New("hash collision occurred")
	ErrorMalformedURL  = errors.New("malformed url input")
	ErrorTryAgain      = errors.New("hash collision, try again with a next base value")
)

// GenerateID 는 url 을 입력받아 hash 된 결과와, 결과가 만들어진 base 와 error 를 return 한다.
// hash 는 sha256 결과(len:64)의 0~9 index 의 값을 base62 encode 한 값이다.
// base 는 hash 생성에 사용한 10 개의 값을 가져온 시작 위치이다, hash 충돌이 발생했을때, 그 다음 위치(base+1) 에서 시도하기 위함이다.
// error 는 잘못된 rawURL 이 입력되거나, hash collision 이 '심각하게 발생((64-10-1) 번의 시도를 연속으로 실패했을때)' 했을때 발생 한다.
// 만약 hash collision 이 발생하거나 0번 index 에 해당하는 값이 0 이라면, base 를 오른쪽으로 1칸씩 shift 를 한다.
func GenerateID(rawURL string, base int) (string, int, error) {
	if _, err := url.Parse(rawURL); err != nil {
		return "", 0, ErrorMalformedURL
	}

	if base >= 64-10 {
		return "", 0, ErrorHashCollision
	}

	return generateID(rawURL, base)
}

func generateID(url string, base int) (string, int, error) {
	sum := sha256.Sum256([]byte(url))

	nArr := make([]int, 64)
	for i := 0; i < len(sum); i++ {
		nArr[2*i], nArr[2*i+1] = int(sum[i])/16, int(sum[i])%16
	}

	tot := 0
	if nArr[base] == 0 {
		return "", base, ErrorTryAgain
	}

	for i := 0; i < len(pow16); i++ {
		tot += nArr[base+i] * pow16[i]
	}

	x := make([]byte, 0)
	for i := 0; i < len(pow62); i++ {
		x = append(x, digits[tot/pow62[i]])
		tot %= pow62[i]
	}

	return string(x), base, nil
}
