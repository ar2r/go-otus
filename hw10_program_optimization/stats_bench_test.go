//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

// Запустил на коде из master ветки
// go test -bench=BenchmarkGetDomainStat -benchmem  -count=10 > old.txt
// Запустил на оптимизированном коде
// go test -bench=BenchmarkGetDomainStat -benchmem  -count=10 > new.txt
// Результат оптимизации
// benchstat old.txt new.txt
func BenchmarkGetDomainStat(b *testing.B) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":6,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@почта.рф","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`
	for i := 0; i < b.N; i++ {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(b, err)
		require.Equal(b, DomainStat{"browsedrive.gov": 1}, result)
	}
}
