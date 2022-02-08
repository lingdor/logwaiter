package fullflag

import "testing"

func TestFileSize(t *testing.T) {

	vals := map[string]int{
		"100m":  1024 * 1024 * 100,
		"100mb": 1024 * 1024 * 100,
		"10g":   1024 * 1024 * 1024 * 10,
		"10gb":  1024 * 1024 * 1024 * 10,
		"10kb":  1024 * 10,
		"10k":   1024 * 10,
		"10b":   10,
		"10":    10,
	}
	for k, v := range vals {
		var result int
		ins := newFileSizeValue(0, &result)
		ins.Set(k)
		if int(*ins) != v {
			t.Errorf("assert %s = %d , want %d", k, int(*ins), v)
		}
	}

}
